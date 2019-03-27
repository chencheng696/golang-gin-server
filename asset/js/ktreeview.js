/*
var tree = [{
	id: "1",
	text: "Node 1", //节点显示的文本值  string
	//icon: "glyphicon glyphicon-play-circle", //节点上显示的图标，支持bootstrap的图标  string
	selectedIcon: "glyphicon glyphicon-ok", //节点被选中时显示的图标       string
	//color: "#ff0000", //节点的前景色      string
	//backColor: "#1606ec", //节点的背景色      string
	nodes: [{
			id: "3",
			text: "Child 1",
			state: { //描述节点的初始状态    Object
				checked: true, //是否选中节点
				//disabled: true, //是否禁用节点
				expanded: true, //是否展开节点
				selected: true //是否选中节点
			},
			nodes: [{
				id: "4",
				text: "Grandchild 1"
			}, {
				id: "5",
				text: "Grandchild 2"
			}]
		}, {
			id: "6",
			text: "Child 2"
		}]
   }, {
	id: "2",
       	text: "Parent 2",
       nodes: [{
		id: "7",
           text: "Child 2",
           nodes: [{
				id: "8",
                text: "Grandchild 3"
            }, {
				id: "9",
                text: "Grandchild 4"
           }]
       }, {
		id: "10",
           text: "Child 2"
       }]
   }, {
	id: "11",
       text: "Parent 3"
   }, {
	id: "12",
       text: "Parent 4"
   }, {
	id: "13",
       text: "Parent 5"
   }];*/
$.fn.extend({
	ktreeview: function(options, callback) {
		var obj = new Object;
		obj.options = options;
		
		var treeid = 'treeview-' + (new Date().getTime());
		var btnid = 'btn-' + (new Date().getTime());
		
		$(document).click(function(e){
			if (treeid != '') {
				if (!$(e.target).parent().hasClass('node-' + treeid))
					$('#' + treeid).hide();
			}
		});
		
		obj.txtno = $(this);
		
		$(obj.txtno).hide();
		
		if ($(obj.txtno).next().prop('tagName') == 'BUTTON' && 
			$(obj.txtno).next().hasClass('ktreeview')) {
			$(obj.txtno).next(':first').remove();
		}
		if ($(obj.txtno).next().prop('tagName') == 'DIV' && 
			$(obj.txtno).next().hasClass('ktreeview')) {
			$(obj.txtno).next(':first').remove();
		}
		
		$(obj.txtno).after('<button id="' + btnid + '" class="btn btn-white btn-default col-xs-10 col-sm-12 ktreeview" style="border-color: #D5D5D5;z-index: 1; margin-bottom: 9px">请选择</button>');
		
		obj.button = $('#' + btnid);
		$(obj.button).after('<div id="' + treeid + '" class="ktreeview" style="display: none;"></div>');
		
		obj.tree = $('#' + treeid);
	
		var defaults = {
			bootstrap2 : false,
			showTags : true,
			//levels : 5,
			multiSelect: false,
			showCheckbox : false,
			checkedIcon : "glyphicon glyphicon-check",
			color: "#000000",
			backColor: "#FFFFFF",
			onNodeSelected : function(event, data) {
				$(obj.tree).hide();
				$(obj.txtno).val(data.id);
				$(obj.button).text(data.text);
				if (callback != null) {
					callback(data);
				}
			}
		};
		var opts = $.extend(defaults, obj.options);
		$(obj.tree).treeview(opts);

		$(obj.button).click(function(event) {
			
			//取消事件冒泡
			var e = arguments.callee.caller.arguments[0] || event;
			if (e && e.stopPropagation)
				e.stopPropagation();
			else if (window.event)
				window.event.cancelBubble = true;
				
			$(obj.tree).show();
			return false;
		});
		
		obj.getName = function (){
			return $(obj.button).text();
		}
		obj.getValue = function (){
			return $(obj.txtno).val();
		}
		obj.setValue = function (v){
			
			var index = 0;
			var nodeId = getTreeNodeId(obj.options.data, v);
			if (nodeId < 0) {
				var arr = obj.tree.treeview('getSelected');
				for (var k in arr) {
					obj.tree.treeview('unselectNode', arr[k].nodeId);
				}
				$(obj.txtno).val('');
				$(obj.button).text('请选择');
				return;
			}
			obj.tree.treeview('checkNode', nodeId);
			
			function getTreeNodeId(data, id) {
				for (var i = 0; i < data.length; i++) {
					data[i].index = index;
					index++;
					
					if (data[i].id == id) {
						return data[i].index;
					} else {
						if (data[i].hasOwnProperty('nodes')) {
							var ret = getTreeNodeId(data[i].nodes, id, index);
							if (ret >= 0) {
								return ret;
							}
						}
					}
				}
				return -1;
			}
		}
		return obj;
	}
});