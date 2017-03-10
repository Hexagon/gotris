define([ 'util/castrato'], function (bus) {

    
	var 
		exports = {},
		socket,
		ws,
		
		url = function (s) {
			
            var l = window.location,
            	host = l.host,
            	dir = l.pathname.substring(0, l.pathname.lastIndexOf('/'));
           
            return ((l.protocol === "https:") ? "wss://" : "ws://") + l.host + dir + s;
            
        };

	exports.connect = function (uri) {

        var ws = new WebSocket(url(uri));
        
		// In
		ws.onmessage    = (m)   => bus.emit('network:message', JSON.parse(m.data));
		ws.onopen   	= ()    => bus.emit('network:connect', ws);
		ws.onclose      = ()    => bus.emit('network:disconnect');

		// Out
		bus.on("controls:change", (data) => ws.send( JSON.stringify({ packet: "key", data: data }) ) );
		bus.on("player:ready", (nickname) => ws.send(  JSON.stringify({ packet: "ready", nickname: nickname }) ) );
		
	};
	
	return exports;
});