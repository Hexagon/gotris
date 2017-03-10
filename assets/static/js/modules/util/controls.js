define(['util/castrato'], function (bus) {

	var 

		// Which keys should trigger stuff
		listenKeys = ["ArrowUp", "ArrowLeft", "ArrowRight", "ArrowDown", "KeyA", "KeyD", "Space"],

		setState = function (prop, val, isEvent) {
			bus.emit("controls:change", { key: prop, state: val });
		};

	/* Keyboard events */
	window.addEventListener('keydown', function (e) {
		
		// Ignore if not in list of interestinhg keys
		if(listenKeys.indexOf(e.code) === -1) return;
		
		// Drop
		if(e.code == "Space") setState("drop", true);
		
		// Rotation & Flip
		if(e.code == "ArrowUp") setState("rotCW", true);
		if(e.code == "KeyA") setState("rotCCW", true);
		if(e.code == "KeyD") setState("rotCW", true);
		
		// Movement
		if(e.code == "ArrowLeft") setState("left", true);
		if(e.code == "ArrowRight") setState("right", true);
		if(e.code == "ArrowDown") setState("down", true);
		
	});
	
	window.addEventListener('keyup', function (e) {
		
		// Ignore if not in list of interestinhg keys
		if(listenKeys.indexOf(e.code) === -1) return;
		
		// Drop
		if(e.code == "Space") setState("drop", false);
		
		// Rotation & Flip
		if(e.code == "ArrowUp") setState("rotCW", false);
		if(e.code == "KeyA") setState("rotCCW", false);
		if(e.code == "KeyD") setState("rotCW", false);
		
		// Movement
		if(e.code == "ArrowLeft") setState("left", false);
		if(e.code == "ArrowRight") setState("right", false);
		if(e.code == "ArrowDown") setState("down", false);

	});

});