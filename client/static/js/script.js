//----- socket.io -----//
var socket = io('', {transports: ['websocket']});//this domain
socket.on('content',function(data){
	editor.getDoc().setValue(data)
});

//----- li active -----//
$('#file-list li').click(function(e) {
	e.preventDefault();
	$('.active').removeClass('active');
	$(this).addClass('active');
});

//----- require content -----//
socket.emit("require", { title: $("#title").text() }, function(){
	console.log("require content of "+ $("#title").text())
})
//----- listen key -----//
document.onkeydown = function(e) {
    switch (e.keyCode) {
		//----- F6 -----//
    case 115:
			break;
		//----- Enter -----//
		case 13:
			var content = editor.getDoc().getValue()
			socket.emit('save', { title: $("#title").text(), content: content });
			break;
    }
};
/*** on ready ***///$(document).ready(function() {});
