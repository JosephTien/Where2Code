var editor = CodeMirror.fromTextArea(document.getElementById("editor"), {
	lineNumbers: true,
	mode: "go",
	//fullScreen: true,
	matchBrackets: true,
	extraKeys: {
		"F11": function(cm) {
			if(cm.getOption("fullScreen")){
				cm.setOption("fullScreen",false);
				$(".navbar").show();
				$(".sidebar").show();
			}else{
				cm.setOption("fullScreen",true);
				 $(".navbar").hide();
				 $(".sidebar").hide();
			}
		},
		"Ctrl-0": function(cm) {
			if(cm.getOption("theme")!="night"){
				cm.setOption("theme","night");
			}else{
				cm.setOption("theme","");
			}
		},
		"Ctrl-+": function(cm) {
			var el = document.getElementsByClassName("CodeMirror")[0];
			var style = window.getComputedStyle(el, null).getPropertyValue('font-size');
			var fontSize = parseFloat(style); 
			el.style.fontSize = (fontSize + 4) + 'px';
		},
		"Ctrl-=": function(cm) {
			var el = document.getElementsByClassName("CodeMirror")[0];
			var style = window.getComputedStyle(el, null).getPropertyValue('font-size');
			var fontSize = parseFloat(style); 
			el.style.fontSize = (fontSize + 4) + 'px';
		},
		"Ctrl--": function(cm) {
			var el = document.getElementsByClassName("CodeMirror")[0];
			var style = window.getComputedStyle(el, null).getPropertyValue('font-size');
			var fontSize = parseFloat(style); 
			el.style.fontSize = (fontSize - 4) + 'px';
		}
	}
});