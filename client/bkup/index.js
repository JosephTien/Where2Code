var box_handler = document.querySelector('.box-handler');
var box_wrapper = box_handler.closest('.box-wrapper');
var boxA = box_wrapper.querySelector('.box');
var isHandlerDragging = false;

document.addEventListener('mousedown', function(e) {
	// If mousedown event is fired from .handler, toggle flag to true
	if (e.target === box_handler) {
		isHandlerDragging = true;
	}
});

document.addEventListener('mousemove', function(e) {
  // Don't do anything if dragging flag is false
	if (!isHandlerDragging) {
		return false;
	}
	// Get offset
	var containerOffsetLeft = box_wrapper.offsetLeft;
	// Get x-coordinate of pointer relative to container
	var pointerRelativeXpos = e.clientX - containerOffsetLeft;
	// Resize box A
	// * 8px is the left/right spacing between .handler and its inner pseudo-element
	// * Set flex-grow to 0 to prevent it from growing
	boxA.style.width = (pointerRelativeXpos - 8) + 'px';
	boxA.style.flexGrow = 0;
});

document.addEventListener('mouseup', function(e) {
	// Turn off dragging flag when user mouse is up
	isHandlerDragging = false;
});