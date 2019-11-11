# -*- coding: utf-8 -*-
"""
/***************************************************************************
 Gisquick plugin
 Publish your projects into Gisquick application
 ***************************************************************************/
"""
import os
import re
import sys
import urllib
import platform
import configparser
from urllib.parse import parse_qs

# Import the PyQt and QGIS libraries
import PyQt5.uic
from qgis.core import Qgis, QgsMapLayer, QgsProject, QgsLayerTreeLayer, QgsLayoutItemLabel
from qgis.PyQt.QtWidgets import QAction, QMessageBox
from qgis.PyQt.QtGui import QIcon
from qgis.PyQt.QtCore import QSettings, QTranslator, qVersion, QCoreApplication

# Initialize Qt resources from file resources.py
from . import resources_rc

from .utils import scales_to_resolutions, resolutions_to_scales, to_decimal_array
from .gisquick_ws import gisquick_ws, WsError


__metadata__ = configparser.ConfigParser()
__metadata__.read(os.path.join(os.path.dirname(__file__), 'metadata.txt'))

class Node(object):
    """
    Tree node element for holding information about layers
    organization in the tree structure.
    """
    name = None
    layer = None
    parent = None
    children = None

    """
    Args:
        name (str): name of the node
        children (List[webgisplugin.Node], List[str]): array of node's children
        layer (qgis.core.QgsMapLayer): qgis layer attached to this node
    """
    def __init__(self, name, children=None, layer=None):
        self.name = name
        self.layer = layer
        self.children = []
        if children:
            self.append(*children)

    """
    Args:
        *nodes (str, webgisplugin.Node): nodes to be appended into this node.
            Strings will be automatically converted to webgisplugin.Node objects.
    """
    def append(self, *nodes):
        for node in nodes:
            if node is not None:
                if not isinstance(node, Node):
                    node = Node(node)
                node.parent = self
                self.children.append(node)

    """
    Search for node with given name, starting with this node and recursively
    continue with it's descendant nodes.

    Args:
        name (str): recursively (from this node) search for a node with given name

    Returns:
        webgisplugin.Node: node with given name, or None if there is no match
    """
    def find(self, name):
        if name == self.name:
            return self
        for child in self.children:
            res = child.find(name)
            if res:
                return res

    """
    Traverse the tree from this node (Post-order) and execute given
    function on each node.

    Args:
        fn (Callback[webgisplugin.Node]): function to be executed
    """
    def cascade(self, fn):
        for child in self.children:
            child.cascade(fn)
        fn(self)


# from qgis.PyQt import QtCore
from PyQt5 import QtCore
from PyQt5.QtCore import QThread, QTimer

class WebGisPlugin(object):

    dialog = None
    project = None
    ws = None

    def __init__(self, iface):
        # Save reference to the QGIS interface
        self.iface = iface
        # initialize plugin directory
        self.plugin_dir = os.path.dirname(__file__)
        # initialize locale
        locale = QSettings().value("locale/userLocale")[0:2]
        localePath = os.path.join(self.plugin_dir, 'i18n', 'gisquick_{}.qm'.format(locale))

        if os.path.exists(localePath):
            self.translator = QTranslator()
            self.translator.load(localePath)

            if qVersion() > '4.3.3':
                QCoreApplication.installTranslator(self.translator)

    def initGui(self):
        # Create action that will start plugin configuration
        self.action = QAction(
            QIcon(":/plugins/gisquick2/img/icon.svg"),
            u"Publish in Gisquick", self.iface.mainWindow())
        self.action.setCheckable(True)
        # connect the action to the run method
        self.action.triggered.connect(self.toggle_tool)

        self.settings_action = QAction(
            QIcon(":/plugins/gisquick2/img/settings.svg"),
            u"Configure", self.iface.mainWindow())
        # connect the action to the run method
        self.settings_action.triggered.connect(self.show_settings)

        # Add toolbar button and menu item
        # self.iface.addToolBarIcon(self.action)
        self.iface.addWebToolBarIcon(self.action)
        self.iface.addPluginToWebMenu(u"&Gisquick2", self.action)
        self.iface.addPluginToWebMenu(u"&Gisquick2", self.settings_action)

    def unload(self):
        if self.ws:
            gisquick_ws.stop()
            self.ws = None

        # Remove the plugin menu item and icon
        self.iface.removePluginWebMenu(u"&Gisquick2", self.action)
        self.iface.removePluginWebMenu(u"&Gisquick2", self.settings_action)
        self.iface.removeWebToolBarIcon(self.action)

    def is_overlay_layer_for_publish(self, layer):
        """Checks whether layer can be published as an overlay layer.

        Args:
            layer (qgis.core.QgsMapLayer): project layer
 
        Returns:
            bool: True if a layer can be published as an overlay layer
        """
        return (layer.type() == QgsMapLayer.VectorLayer or
            (layer.type() == QgsMapLayer.RasterLayer and layer.providerType() != "wms"))

    def is_base_layer_for_publish(self, layer):
        """Checks whether layer could be published as a base layer.

        Args:
            layer (qgis.core.QgsMapLayer): project layer
 
        Returns:
            bool: True if a layer can be published as a base layer
        """
        return layer.type() == QgsMapLayer.RasterLayer and layer.providerType() == "wms"

    def map_units(self):
        """Returns units name of the project (map).

        Returns:
            str: map units name ('meters', 'feet', 'degrees', 'miles' or 'unknown')
        """
        return {
            0: 'meters',
            1: 'feet',
            2: 'degrees',
            3: 'unknown',
            7: 'miles'
        }[self.iface.mapCanvas().mapUnits()]

    def scales_to_resolutions(self, scales):
        """Converts map scales to tile resolutions (with fixed DPI=96).

        Args:
            scales (List[int]): array of map scales

        Returns:
            List[Decimal]: array of computed tile resolutions
        """
        return scales_to_resolutions(scales, self.map_units())

    def resolutions_to_scales(self, resolutions):
        """Converts tile resolutions to map scales (with fixed DPI=96).

        Args:
            resolutions (List[Decimal]): array of tile resolutions

        Returns:
            List[int]: array of computed map scales
        """
        return resolutions_to_scales(resolutions, self.map_units())

    def filter_visible_resolutions(self, resolutions, layer):
        """Filters given tile resolutions by layer's visibility settings.

        Args:
            resolutions (List[Decimal]): array of tile resolutions
            layer (qgis.core.QgsMapLayer): map layer

        Returns:
            List[Decimal]: array of visible tile resolutions
        """
        if layer.hasScaleBasedVisibility():
            max_scale_exclusive = layer.maximumScale()
            min_scale_inclusive = layer.minimumScale()
            max_res_exclusive, min_res_inclusive = self.scales_to_resolutions(
                [max_scale_exclusive, min_scale_inclusive]
            )
            return [res for res in resolutions if res >= min_res_inclusive and res < max_res_exclusive]
        return resolutions

    def wmsc_layer_resolutions(self, layer):
        """Returns visible resolutions of given WMSC layer.

        Args:
            layer (qgis.core.QgsRasterLayer): raster layer (WMSC)

        Returns:
            List[Decimal]: array of layer's visible tile resolutions
        """
        layer_resolutions = layer.dataProvider().property('resolutions')
        if layer_resolutions:
            layer_resolutions = to_decimal_array(layer_resolutions)
            if layer.hasScaleBasedVisibility():
                layer_resolutions = self.filter_visible_resolutions(layer_resolutions, layer)
            if layer_resolutions:
                return sorted(layer_resolutions, reverse=True)
            return []
        return None

    def project_layers_resolutions(self):
        """Returns list of possible tile resolutions for current project.

        Returns:
            List[Decimal]: project tile resolutions
        """
        # compute resolutions as an union of resolutions calculated from project's
        # map scales and resolutions of all WMSC layers.
        project_tile_resolutions = set()

        # collect set of all resolutions from WMSC base layers
        base_layers = {
            layer.id(): layer
            for layer in QgsProject.instance().mapLayers().values()
                if self.is_base_layer_for_publish(layer)
        }
        for layer in list(base_layers.values()):
            layer_resolutions = self.wmsc_layer_resolutions(layer)
            if layer_resolutions:
                project_tile_resolutions.update(layer_resolutions)

        wmsc_layers_scales = self.resolutions_to_scales(project_tile_resolutions)
        scales, ok = self.project.readListEntry("Scales", "/ScalesList")
        if ok and scales:
            scales = [int(scale.split(":")[-1]) for scale in scales]
            # filter duplicit scales
            scales = [scale for scale in scales if scale not in wmsc_layers_scales]
            project_tile_resolutions.update(
                self.scales_to_resolutions(sorted(scales, reverse=True))
            )

        project_tile_resolutions = sorted(project_tile_resolutions, reverse=True)
        return project_tile_resolutions

    def layers_list(self):
        """Returns array of all project's layers.

        Returns:
            List[qgis.core.QgsMapLayer]: project's layers
        """
        # legend_iface = self.iface.legendInterface().layers()
        return [tree_layer.layer() for tree_layer in QgsProject.instance().layerTreeRoot().findLayers()]


    def get_layer_attributes(self, layer):
        fields = layer.fields()
        attributes_data = []
        excluded_attributes = layer.excludeAttributesWfs()
        conversion_types = {
            'BIGINT': 'INTEGER',
            'INTEGER64': 'INTEGER',
            'REAL': 'DOUBLE',
            'STRING': 'TEXT',
            'INT2': 'INTEGER',
            'INT4': 'INTEGER',
            'INT8': 'INTEGER',
            'NUMERIC': 'DOUBLE',
            'FLOAT8': 'DOUBLE',
            'VARCHAR': 'TEXT',
            'CHARACTER': 'TEXT'
        }
        for field in fields:
            if field.name() in excluded_attributes:
                continue
            field_type = field.typeName().upper()
            if field_type in conversion_types:
                field_type = conversion_types[field_type]
            attribute_data = {
                'name': field.name(),
                'type': field_type,
                #'length': field.length(),
                #'precision': field.precision()
            }
            if field.comment():
                attribute_data['comment'] = field.comment()
            alias = layer.attributeAlias(fields.indexFromName(field.name()))
            if alias:
                attribute_data['alias'] = alias
            attributes_data.append(attribute_data)

        return attributes_data

    def get_project_layers(self, skip_layers_with_error=False):
        dbname_pattern = re.compile("dbname='([^']+)'")
        project = QgsProject.instance()
        project_dir = project.absolutePath() + os.path.sep

        non_identifiable_layers = project.readListEntry("Identify", "/disabledLayers")[0] or []
        wfs_layers = project.readListEntry("WFSLayers", "/")[0] or []
        map_canvas = self.iface.mapCanvas()
        map_settings = map_canvas.mapSettings()

        if project.layerTreeRoot().hasCustomLayerOrder():
            layers_order = project.layerTreeRoot().customLayerOrder()
        else:
            layers_order = self.layers_list()

        def visit_node(tree_node):
            if isinstance(tree_node, QgsLayerTreeLayer):
                layer = tree_node.layer()
                layer_type = {
                    0: "vector",
                    1: "raster",
                    # 2: PluginLayer
                    # 3: MeshLayer
                }[layer.type()]
                source = layer.source()
                provider_type = layer.providerType()
                uri = ""

                if provider_type == "wms":
                    source_params = parse_qs(source)
                    uri = source_params["url"][0]
                elif provider_type == "postgres":
                    dp = layer.dataProvider()
                    uri = "postgresql://%s:%s" % (dp.uri().host(), dp.uri().port())
                elif provider_type in ("ogr", "gdal"):
                    uri = "file://%s" % source.split("|")[0]
                elif provider_type == "spatialite":
                    match = dbname_pattern.search(source)
                    if match:
                        uri = "file://%s" % match.group(1)
                else:
                    uri = source

                extent = layer.extent()
                if not extent.isEmpty():
                    extent = map_settings.layerExtentToOutputExtent(
                        layer,
                        layer.extent()
                    ).toRectF().getCoords()
                else:
                    extent = None
                info = {
                    "name": layer.name(),
                    "serverName": layer.shortName() if hasattr(layer, "shortName") else layer.name(),
                    "wfs": layer.id() in wfs_layers,
                    "provider_type": provider_type,
                    "projection": layer.crs().authid(),
                    "type": layer_type,
                    "source": uri,
                    "extent": extent,
                    "visible": project.layerTreeRoot().findLayer(layer.id()).isVisible(),
                    "metadata": {
                        "title": layer.title(),
                        "abstract": layer.abstract(),
                        "keyword_list": layer.keywordList()
                    }
                }
                # if layer.isSpatial()
                if layer_type == "vector":
                    info["geom_type"] = ('POINT', 'LINE', 'POLYGON', None, None)[layer.geometryType()]
                    info["labels"] = layer.labelsEnabled()
                    info["attributes"] = self.get_layer_attributes(layer)
                    if info["attributes"]:
                        fields = layer.fields()
                        info['pk_attributes'] = [
                            fields.at(index).name() for index in layer.dataProvider().pkAttributeIndexes()
                        ]

                if provider_type == "wms":
                    info["format"] = source_params["format"][0]
                    info["url"] = source_params["url"][0]
                    info["dpi"] = layer.dataProvider().dpi()
                    if "layers" in source_params:
                        info["wms_layers"] = source_params["layers"]

                if layer in layers_order:
                    info["drawing_order"] = layers_order.index(layer)
                if layer.attribution():
                    info["attribution"] = {
                        "title": layer.attribution(),
                        "url": layer.attributionUrl()
                    }
                return info
            else:
                children = []
                for child_tree_node in tree_node.children():
                    try:
                        info = visit_node(child_tree_node)
                        if info:
                            children.append(info)
                    except Exception as e:
                        if not skip_layers_with_error:
                            msg = "Failed to gather info from layer: '%s'" % child_tree_node.name()
                            raise WsError(msg, 405) from e
                return {
                    "name": tree_node.name(),
                    "layers": children
                }

        root_node = self.iface.layerTreeView().layerTreeModel().rootGroup()
        return visit_node(root_node)["layers"]

    def get_print_templates(self):
        composer_templates = []
        project_layout_manager = QgsProject.instance().layoutManager()
        for layout in project_layout_manager.layouts():
            map = layout.referenceMap()
            units_conversion = map.mapUnitsToLayoutUnits()
            composer_data = {
                'name': layout.name(),
                'width': layout.layoutBounds().width(),
                'height': layout.layoutBounds().height(),
                'map': {
                    'name': 'map0',
                    'x': map.pagePos().x(),
                    'y': map.pagePos().y(),
                    'width': map.extent().width() * units_conversion,
                    'height': map.extent().height() * units_conversion
                },
                'labels': [
                    item.id() for item in list(layout.items())
                        if isinstance(item, QgsLayoutItemLabel) and item.id()
                ]
            }
            grid = map.grid()
            if grid.enabled():
                composer_data['map']['grid'] = {
                    'intervalX': grid.intervalX(),
                    'intervalY': grid.intervalY(),
                }
            composer_templates.append(composer_data)
        return composer_templates

    def get_project_info(self, skip_layers_with_error=False):
        project = QgsProject.instance()
        project_crs = project.crs()
        map_canvas = self.iface.mapCanvas()

        scales, _ = project.readListEntry("Scales", "/ScalesList")
        scales = [int(s.split(":")[1]) for s in scales]

        data = {
            "file": project.absoluteFilePath(),
            "directory": project.absolutePath(),
            "title": project.title() or project.readEntry("WMSServiceTitle", "/")[0],
            "layers": self.get_project_layers(skip_layers_with_error),
            "composer_templates": self.get_print_templates(),
            "projection": {
                "code": project_crs.authid(),
                "is_geographic": project_crs.isGeographic(),
                "proj4": project_crs.toProj4()
            },
            "units": self.map_units(),
            "scales": scales,
            "position_precision": {
                "automatic": project.readBoolEntry("PositionPrecision", "/Automatic")[0],
                "decimal_places": project.readNumEntry("PositionPrecision", "/DecimalPlaces")[0]
            },
            "extent": map_canvas.fullExtent().toRectF().getCoords(),
            "server": {
                "wms_add_geometry": project.readBoolEntry ("WMSAddWktGeometry", "")[0]
            }
        }
        return data

    def get_settings(self):
        return QSettings(QSettings.IniFormat, QSettings.UserScope, "Gisquick", "gisquick2")

    def show_settings(self):
        settings = self.get_settings()
        dialog_filename = os.path.join(self.plugin_dir, "ui", "settings.ui")
        dialog = PyQt5.uic.loadUi(dialog_filename)
        dialog.server_url.setText(settings.value("server_url", ""))
        dialog.username.setText(settings.value("username", ""))
        dialog.password.setText(settings.value("password", ""))

        dialog.show()
        res = dialog.exec_()
        if res == 1:
            settings.setValue("server_url", dialog.server_url.text())
            settings.setValue("username", dialog.username.text())
            settings.setValue("password", dialog.password.text())


    def on_project_change(self, *args):
        gisquick_ws.send("ProjectChanged")

    def on_project_closed(self, *args):
        def debounced():
            # filter events caused by switching between projects
            if not QgsProject.instance().absoluteFilePath():
                gisquick_ws.send("ProjectChanged")

        QTimer.singleShot(300, debounced)

    def toggle_tool(self, active):
        """Display dialog window for publishing current project.

        During a configuration process (wizard setup), plugin will hold actual metadata
        object in 'WebGisPlugin.metadata' property. If metadata from previous publishing
        still exist, they will be loaded and stored in 'WebGisPlugin.last_metadata' property.
        """

        def callback(msg):
            msg_type = msg["type"]
            data = msg.get("data")
            if msg_type == "ProjectInfo":
                if QgsProject.instance().absolutePath():
                    if data:
                        skip_layers_with_error = data.get("skip_layers_with_error", False)
                    return self.get_project_info(skip_layers_with_error=skip_layers_with_error)
                raise WsError("Project is not opened", 404)

            elif msg_type == "ProjectDirectory":
                dir_path = QgsProject.instance().absolutePath()
                if dir_path:
                    return dir_path
                raise WsError("Project is not opened", 404)
            else:
                raise ValueError("Unknown message type: %s" % msg_type)


        project = QgsProject.instance()
        if active:
            settings = self.get_settings()
            server_url = settings.value("server_url")
            username = settings.value("username")
            password = settings.value("password")
            if not server_url or not username or not password:
                self.show_settings()
                server_url = settings.value("server_url")
                username = settings.value("username")
                password = settings.value("password")

            plugin_ver = __metadata__["general"].get("version")
            client_info = "GisquickPlugin/%s (%s %s; QGIS %s)" % (plugin_ver, platform.system(), platform.machine(), Qgis.QGIS_VERSION)

            class WebsocketServer(QThread):
                finished = QtCore.pyqtSignal(int)

                def run(self):
                    # print("Starting WS", "server:", server_url, "user:", username)
                    res = gisquick_ws.start(server_url, username, password, client_info, callback)
                    self.finished.emit(res)

            def on_finished(res):
                self.ws = None
                if self.action.isChecked():
                    self.action.setChecked(False)
                if res != 0:
                    QMessageBox.warning(None, 'Warning', 'Failed to connect!')
            self.ws = WebsocketServer()
            self.ws.finished.connect(on_finished)
            self.ws.start()

            # project.isDirtyChanged.connect(self.on_project_change)
            project.readProject.connect(self.on_project_change)
            project.projectSaved.connect(self.on_project_change)
            project.cleared.connect(self.on_project_closed)
        else:
            gisquick_ws.stop()
            project.readProject.disconnect(self.on_project_change)
            project.projectSaved.disconnect(self.on_project_change)
            project.cleared.disconnect(self.on_project_closed)
