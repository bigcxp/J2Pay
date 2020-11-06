const API_HOST = '/';


function postform(form,url,datas,fn,isjson,method){
	var data,$form;
	if(!form){form=$('form').get(0);}
	$form=$(form);
	if(!datas){
		data = (true==isjson) ? JSON.stringify(form2json(form)) : trim_data($form.serialize());
	}else{
		if(typeof(data)=='object'){data=JSON.stringify(data);}
	}
	if(!url && !!form['action']){url=form.action;}
	if(!url){url=location.href;}
	$.ajax({
		'crossDomain':true,
		'xhrFields':{'withCredentials':true},
		'processData':(true==isjson)?false:true,
		'url':url,
		'headers':{'Authorization':sessionStorage.getItem('token')},
		'type':method||'POST',
		'data':data,
		'contentType': (true==isjson)?'application/json':'application/x-www-form-urlencoded',
		success:function(res){
			if(typeof(fn)=='function'){fn(res);return;}
			if(res.code=='1'){
				alertok(res.msg);
			}else{
				alerterr(res.msg);
			}
		}
	});
}
function postjson(fn,form,url,datas){
	postform(form,url,datas,fn,true);
}

function msgok(str){layer.msg(str,{icon:1});}
function msgerr(str){layer.msg(str,{icon:2,shade:0.3,time:2000});}
function tips(el,str){
	layer.tips(str,el,{tips:[3,'#388E3C']});
}
function tipserr(el,str){
	layer.tips(str,el,{
		tips:[3,'#F44336'],
		time:3300
	});
	var dom=$(el);
	dom.addClass('errbg').focus();
	setTimeout(function(){dom.removeClass('errbg')},3000);
}
function alertok(str,fn){
	layer.alert(str,{icon:1},
		function(){if(typeof(fn)=='function'){fn()}}
	);
}
function alerterr(str){layer.alert(str,{icon:2,anim:6});}
function objlen(obj){
	var len=0;
	$.each(obj,function(){len++;});
	return len;
}
function trim_data(data){
	var res={};
	if(typeof(data)=='string'){
		return data.replace(/(%20)+\&/g,'&').replace(/\=(%20)+/g,'=')	//space
			.replace(/(%0A)+\&/g,'&').replace(/\=(%0A)+/g,'=')	// \n
			.replace(/(%09)+\&/g,'&').replace(/\=(%09)+/g,'=')	// \t
	}
	$.each(data,function(k,v){
		res[k]=$.trim(v);
	})
	return $.param(res);
}
function form2json(el){
	if(!el){el='form:eq(0)';}
	var o={};
	var a=$(el).serializeArray();
	$.each(a, function(i,v) {
		if (o[v.name]) {
			if (!o[v.name].push) {
				o[v.name] = [o[v.name]];
			}
			o[v.name].push(v.value || '');
		} else {
			o[v.name] = v.value || '';
		}
	});
	return o;
}

function open_dlg(title,url,size,fn){
	if(!url){return;}
	var idx=layer.open({
		type:2,
		title:title||'Dialog',
		content:url,
		area: size||['90%','90%'],
		btn:['close'],
		anim:-1,
		isOutAnim:false,
		yes:function(){layer.close(idx)},
		success:function(){
			{if(typeof(fn)=='function'){fn()}}
		}
	});
}
function close_dlg(){
	try{
		window.parent.layer.closeAll('iframe')
	}catch(e){}
}

function confirmdel(url,data) {
	var idx = layer.prompt({
		value: '',
		title: '是否确认删除',
		success: function () {
			$('.layui-layer-content input').attr('maxlength', '6').attr('type', 'phone').attr('placeholder','请输入6位Google验证码')
		}
	}, function (value, index, elem) {
		if ($.trim(value) == '') {
			return;
		}
		var loadingidx=layer.load(0,{shade:[0.2,'#000']});
		$.ajax({
			type: "DELETE",
			url: url,
			dataType: "json",
			data: data,
			'headers':{'Authorization':sessionStorage.getItem('token')},
			success: function (res){
				window.top.layer.alert(res.msg,{icon: (res.code == "1") ? 1 : 2})
				if (res.code == "1") {
					try {
						window['dataTable'].reload();
					} catch (e) {}
					layer.closeAll();
				}
			},
			complete:function(){
				layer.close(loadingidx);
			}
		});
	});
}


function h5post(form,url,method){
	var $form=$(form);
	var btn_submit = $form.find('[type=submit]');
	if(!method){method='POST'}
	btn_submit.attr('disabled', 'disabled').addClass('layui-disabled').append('<i class="layui-icon layui-icon-loading layui-icon layui-anim layui-anim-rotate layui-anim-loop"></i>')
	postform(form,url,null,function(res){
		btn_submit.removeAttr('disabled').removeClass('layui-disabled').find('i.layui-icon').remove()
		if(res.code=='1'){
			alertok(res.msg,function(){
				try{parent['dataTable'].reload();}catch(e){}
				try{parent.layer.closeAll('iframe')}catch(e){}
			});
		}else{
			alerterr(res.msg);
		}
	},false,method)
}

function parseHash(){
	var hash=(location.hash).replace('#','');
	var arr=hash.split('?');
	if(arr.length<2){return {};}
	var obj={}
	var arr1=arr[1].split('&');
	$.each(arr1,function(k,v){
		var ss=v.split('=')
		obj[ss[0]]=ss[1];
	});
	return obj;
}

//格式化是否
function render_bool(v,yes,no){
	return (true==v)? '<font color="green">'+(!!yes?yes:'是 √')+'</font>': '<font color="red">'+(!!no?no:'否 ×')+'</font>';
}
function render_timestamp(t){
	var dt=new Date(parseInt(t)*1000);
	return dt.getFullYear()+'-'+(dt.getMonth()+1)+'-'+dt.getDate()+' '+dt.getHours()+':'+dt.getMinutes()+':'+dt.getSeconds();
}

/*
表单验证
*/

function chk_rate(v, item) {
	if (v == '') {
		return
	}
	v = 1 * v;
	if (1 * v == 0) {
		return
	}
	if (v < 0 || v > 1000) {
		return '错误的金额格式'
	}
}
function chk_username(v,d){
	var reg = /^[a-zA-Z0-9_]{5,16}$/;
	if(!reg.test(v)){
		return '用户名由5~16位 数字/字母/下划线 组成!'
	}
	if(/^[0-9]/.test(v)){
		return '用户名用户名不能以数字开头!'
	}
}
function chk_password(v){
	if(!/^[\S]{8,30}$/.test(v)){
		return '密码长度8~30位!';
	};
	if(!(/[0-9]/.test(v) && /[a-z]/.test(v) && /[A-Z]/.test(v))){
		return '密码必须包含数字和大小写字母';
	};
}
function chk_repassword(v,d){
	var fo=form.val('password');
	if(v!=fo.password){
		return '两次输入的密码不同！';
	}
}