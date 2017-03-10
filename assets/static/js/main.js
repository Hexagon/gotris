require.config({ baseUrl: 'static/js/modules' });

define(['viewport', 'network', 'util/castrato', 'game', 'util/controls'], function(viewport, network, bus, game) {

	var
		
		nickname,
		
		elements = {
			screens: {
				login: document.getElementById('viewMenu'),
				game: document.getElementById('viewGame')
			},
			inputs: {
				nickname: document.getElementById('txtNickname')
			},
			buttons: {
				start: document.getElementById('btnStart')
			},
			containers: {
				highscore: document.getElementById('highscore')
			}
		},
		
		showScreen = function (e) {
			for(var s in elements.screens) elements.screens[s].style.display = 'none';
			e.style.display = 'block';
		},
		
		validateInput = function() {
			var inputBaseClass = elements.inputs.nickname.className.replace(/ invalid/,'');
			nickname = elements.inputs.nickname.value.trim();
			if (nickname.length < 1) {
				elements.inputs.nickname.className = inputBaseClass + ' invalid';
				return false;
			} else {
				elements.inputs.nickname.className = inputBaseClass;
				return true;
			}
		},
		
		startGame = function() {
			if (validateInput()) {
				showScreen(elements.screens.game);
				bus.emit("player:ready", nickname);
			} else {
				elements.inputs.nickname.focus();
			}
		};
    
	// Initialize viewport
	viewport.create();
	
	bus.on("game:updated", function () {
		viewport.redraw();
	})
	
	// Show login screen
	showScreen(elements.screens.login);
	
	// Focus nick input
	elements.inputs.nickname.focus();
	
	// Start game on click of button
	network.connect('/ws');
	elements.inputs.nickname.addEventListener('keyup', validateInput);
	elements.inputs.nickname.addEventListener('keydown', function (e) {
		if (e.code == 'Enter') {
			startGame();
		}
	});
	elements.buttons.start.addEventListener('click', startGame);
	
	// Fetch highscore
	var xhr = new XMLHttpRequest();
	xhr.open('GET', '/api/highscores')
	xhr.onload = function() {
	    if (xhr.status === 200) {
	    	
	        var res = JSON.parse(xhr.responseText),
	        	html = '',
	        	
	        	current = 0,
	        	max = 10;
	        	
			res.highscore.forEach(function(hs) {
				if (current++ < max) {
					html += "<div class=\"highscore-entry\"><h5 class=\"right no-margin\">" + hs.Score + "</h5><h5 class=\"no-margin\">" + hs.Nickname + "</h5></div>"
				}
			});
			
			elements.containers.highscore.innerHTML = html;
			
	    }
	};
	xhr.send();

});