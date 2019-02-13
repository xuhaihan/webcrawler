var util = new Object();
/****** 工具类方法 ---  start ******/

/******根据条件查询实体类对象或list ---  start ******/
/**
 * 根据key查询json对象
 */
util.findObjByKey=function(list,key,val){
    var obj;
    var option={};
    option[key] = val;
    var nl = util.findList(list,option);
    if(nl&&nl.length>0){
        obj = nl[0];
    }
    return obj;
};
/**
 * 通过条件查询获取json list
 * @param ol 原json list
 * @param pram_option 条件【json对象】
 * @returns 新json list
 */
util.findList=function(ol,pram_option){
    var nl=[];
    if(ol&&pram_option){
        for ( var i=0; i<ol.length; i++) {
            var obj = ol[i];
            var flag = util.checkObj(obj,pram_option);
            if(flag){
                nl.push(obj);
            }
        }
    }
    return nl;
};
/**
 * 验证json对象是否符合条件
 * @param obj json对象
 * @param pram_option 条件【json对象】
 * @returns boolean flag
 */
util.checkObj=function(obj,pram_option){
    var flag = true;
    for(var key in pram_option){
        if(obj[key]!=pram_option[key]){
            flag = false;
            break;
        }
    }
    return flag;
};
/******根据条件查询实体类对象或list ---  end ******/

/**
 * 验证是否是数字
 */
util.checkNum=function(z){
    if(!z)
        return true;
    var pattern = new RegExp("^[0-9]*$");
    return pattern.test(z);
};

util.checkString=function(s) {
    if(!s)
        return true;
    //var reg = new RegExp("[`~!@#$^&*()=|{}':;',\\[\\].<>/?~！@#￥……&*（）—|{}【】‘；：”“'。，、？]");
    //var reg = new RegExp("[`~!@#$^&*()=|{}':;',\\[\\]<>/?~！@#￥……&*（）—|{}【】‘；：”“'。，、？]");
    var reg = new RegExp("[`~!@#$^&*=|{}':;',\\[\\]<>/?~！@#￥……&*—|{}【】‘；：”“'。，、？]");
    return reg.test(s);
};

/**
 * 检查是否有空格
 */
util.checkSpace=function(s) {
    var f = true;
    if(s.indexOf(" ")==-1&&s.indexOf("　")==-1){
        f = false;
    }
    return f;
};

/**
 * 去空格
 */
util.trim=function(str)
{
    for(var  i  =  0  ;  i<str.length  &&  str.charAt(i)==" "  ;  i++  )  ;
    for(var  j  =str.length;  j>0  &&  str.charAt(j-1)==" "  ;  j--)  ;
    if(i>j)  return  "";
    return  str.substring(i,j);
};

util.get=function(id){return document.getElementById(id);};

/**
 * 验证cet的准考证号
 */
util.checkCetZkzh=function(z){
    var f = false;
    if(z){
        var t = z.charAt(0);
        if(t!="F"&&t!="S"){
            if(!util.checkNum(z))f = true;
            else{
                var t = z.charAt(9);
                if(isNaN(t))f = true;
            }
        }else{
            if(!util.checkNum(z.substring(1)))f = true;
        }
    }
    return f;
};

util.checkTime = function(startTime){
    var startDate = new Date(startTime);
    var t = startDate.getTime();
    var d=new Date().getTime();
    if(d>=t){
        return "";
    }
    return "对不起，请于"+util.showtime(startDate)+"再来查询！";
};

util.showtime=function(d) {
    var hours = d.getHours();
    var minutes = d.getMinutes();
    var timeValue = d.getFullYear()+"年"+(d.getMonth()+1)+"月"+d.getDate()+"日"+((hours >= 12)?"下午":"上午");
    timeValue += ((hours >12) ? hours -12 :hours);
    timeValue += ((minutes < 10) ? ":0" : ":") + minutes;
    return timeValue;
};

util.nec=function(ca,e,daz,dax,das){
    var p = [];
    p.push("ca="+ca);
    p.push("e="+e);
    if(daz){
        p.push("daz="+daz);
        p.push("dax="+encodeURIComponent(dax));
        p.push("das="+das);
    }
    var h = "http://tj.neea.edu.cn/tj.gif?"+p.join("&")+"&t="+(+new Date);
    util.nec_load(h);
};

util.nec_load=function(a){
    var e = new Image;
    e.onload = e.onerror = e.onabort = function() {
        e.onload = e.onerror = e.onabort = null;
        e = null;
    };
    e.src = a;
};