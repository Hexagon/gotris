define([ 'util/castrato'], function (bus) {
    
	var 
		exports = {
			colors: {
				'I': "rgb(0,255,255)",	// Cyan
				'O': "rgb(255,255,0)",	// Yellow
				'T': "rgb(196,0,196)",	// Purple
				'S': "rgb(0,255,0)",	// Green
				'Z': "rgb(255,0,0)",	// Red
				'J': "rgb(32,32,255)",	// Blue
				'L': "rgb(255,165,0)"	// Orange
			},
			data: false,
			grid: false,
			playing: true
		};
	
	bus.on('network:message', function (o) {
		
		if (o.Position) {
			exports.data = o;
		} else if (o.Data) {
			exports.grid = o;
		} else if (o.gameOver) {
			exports.playing = false;
		} else if (o.ready !== undefined) {
			bus.emit('game:ready', o);
		}
		
		bus.emit('game:updated');
		
	});
	
	return exports;
});