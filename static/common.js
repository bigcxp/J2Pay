const API_HOST = '/';




function postjson(form,url,datas,fn){
	if(!url){url=location.href;}
	if(!data){data=form2json($('form'));}
	if(typeof(data)=='object'){data=JSON.stringify(data);}
	$.ajax({
		'crossDomain':true,
		'xhrFields':{'withCredentials':true},
		'url':url,
		'processData':false,
		'type':'POST',
		'data':data,
		'contentType': 'application/json',
		success:function(res){
			if(typeof(fn)=='function'){fn(res);}
		}
	})
}
function postform(form,url,datas,fn,isjson){
	var data,$form;
	if(!form){form=$('form').get(0);}
	$form=$(form);
	if(!datas){
		data = isjson ? form2json($form) : $form.serialize();
	}else{
		if(typeof(data)=='object'){data=JSON.stringify(data);}
	}
	data=trim_data(data);
	if(!url && !!form['action']){url=form.action;}
	if(!url){url=location.href;}
	$.ajax({
		'crossDomain':true,
		'xhrFields':{'withCredentials':true},
		'url':url,
		'type':'POST',
		'data':data,
		success:function(res){
			if(typeof(fn)=='function'){fn(res);}
			if(res.code=='0'){
				alertok(res.msg);
			}else{
				alerterr(res.msg);
			}
		}
	});
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