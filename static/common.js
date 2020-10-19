const API_HOST = '/';


function postform(form,url,datas,fn,isjson){
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
		'type':'POST',
		'data':data,
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