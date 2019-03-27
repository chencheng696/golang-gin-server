function ajaxPost(url, data, success) {
    $.ajax({
        type: "POST",
        url: url,
        data: data,
        error: function (XMLHttpRequest, textStatus, errorThrown) 
        {
            document.location.reload();
        },
        success: function (obj)
        {
            if (typeof(obj) == 'object' && obj.ret == 9999) {
				showAlertDialog('太长时间没有操作，请重新登录', function (){
					document.location.reload();
				});
                return;
            }
            success(obj);
        }
    });
}

function showAlertDialog(msg, func) {
    bootbox.dialog({
        message: "<span class='bigger-110'>" + msg + "</span>",
        buttons:
        {
            "button" :
            {
                "label" : "OK",
                "className" : "btn-sm"
            }
        },
		callback: function () {
			if (func != null)
				func();
		}
    });
}

function showModalError(err, modal) {
    for (var i in err) {
        var key = '#errmodal-' + i;
        //key = key.replace(/_/g, "-");
        if ($(modal).find(key).length > 0) {
            $(modal).find(key).text(err[i]);
        }
    }
}

function clearModalError(modal) {
    $(modal).find('div').each(function(){
        if (typeof($(this).attr("id")) == 'undefined')
            return;
        if ($(this).attr("id").indexOf('errmodal-') == 0) {
            $(this).html('');
        }
    });
}

function showModalValue(data, modal) {
	var max;
    for (var i in data) {
        var key = '#modal-' + i.toLowerCase();
        if ($(modal).find(key).length > 0) {
			max = $(modal).find(key).attr('maxlength');
			if (typeof(max) != 'undefined' && parseInt(max) > 0) {
				$(modal).find(key).val(data[i].substr(0, max));
			} else {
        			$(modal).find(key).val(data[i]);
			}
        }
    }
}

function clearModalValue(modal) {
    $(modal).find('input').each(function(){
        if (typeof($(this).attr("id")) == 'undefined')
            return;
        if ($(this).attr("id").indexOf('modal-') == 0) {
            $(this).val('');
        }
    });
    $(modal).find('select').each(function(){
        if (typeof($(this).attr("id")) == 'undefined')
            return;
        if ($(this).attr("id").indexOf('modal-') == 0) {
            $(this).val('');
        }
    });
    $(modal).find('textarea').each(function(){
        if (typeof($(this).attr("id")) == 'undefined')
            return;
        if ($(this).attr("id").indexOf('modal-') == 0) {
            $(this).val('');
        }
    });
}

function CheckAjaxReturnData(obj) {
	if (obj == null)
		return '服务器发生未知错误！';
	
	if (typeof(obj) != 'object')
		return '服务器发生未知错误！';
	
	return '';
}

function removeKeyModal(data) {
    var v = {};
    for (var k in data) {
        v[k.replace(/modal-/g, "")] = data[k];
    }
    return v;
}

function submitPageNext() {
    $('#pageNo').val(parseInt($('#pageNo').val()) + 1);
    $('#cmd').val('list_search');
    $("#mainform").submit();
}

function submitPagePrevious() {
    $('#pageNo').val(parseInt($('#pageNo').val()) - 1);
    $('#cmd').val('list_search');
    $("#mainform").submit();
}

function submitPageNo(i) {
    $('#pageNo').val(i);
    $('#cmd').val('list_search');
    $("#mainform").submit();
}

function setFirstInputFocus(element) {
	element.find('input').each(function(){
		if ($(this).is(':visible')) {
			$(this).focus();
			return false;
		}
	});
}