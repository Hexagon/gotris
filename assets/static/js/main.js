require.config({ baseUrl: 'static/js/modules' });

define(['viewport', 'network', 'util/castrato', 'game', 'util/controls'], function(viewport, network, bus, game) {

	var
		
		nickname,
		
		elements = {
			screens: {
				login: document.getElementById('viewMenu'),
				gameSingle: document.getElementById('viewGameSingle'),
				gameBattle: document.getElementById('viewGameBattle'),
			},
			inputs: {
				nickname: document.getElementById('txtNickname')
			},
			buttons: {
				startSingle: document.getElementById('btnStartSingle'),
				startBattle: document.getElementById('btnStartBattle'),
			},
			containers: {
				hsAth: document.getElementById('hsAth'),
				hsWeek: document.getElementById('hsWeek'),
				message: document.getElementById('message')
			}
		},
		
		showScreen = function (e) {
			
			for(var s in elements.screens) elements.screens[s].style.display = 'none';
			e.style.display = 'block';
		},
		
		startGameSingle = function() {
			bus.emit("player:ready", elements.inputs.nickname.value.trim());
		},
		
		startGameBattle = function() {
			bus.emit("player:awaiting", elements.inputs.nickname.value.trim());
		};
    
	// Initialize viewport
	viewport.create('#gameSingle');
	
	bus.on("game:updated", function () {
		viewport.redraw();
	})
	
	bus.on("game:ready", function (m) {
		if (m.ready) {
			showScreen(elements.screens.gameSingle);
			
		} else {
			// Indicate error with input box border
			var inputBaseClass = elements.inputs.nickname.className.replace(/ invalid/,'');
			elements.inputs.nickname.className = inputBaseClass + ' invalid';
			
			// Show error as text
			elements.containers.message.className = elements.containers.message.className.replace(/hidden/,'');
			if (m.error) {
				elements.containers.message.innerHTML = m.error;
			} else {
				elements.containers.message.innerHTML = "Unknown error";
			}
			
		}
	});
	
	// Show login screen
	showScreen(elements.screens.login);
	
	// Focus nick input
	elements.inputs.nickname.focus();
	
	// Start game on click of button
	network.connect('/ws');
	elements.inputs.nickname.addEventListener('keydown', function (e) {
		if (e.code == 'Enter') {
			startGameSingle();
		}
	});
	elements.buttons.startSingle.addEventListener('click', startGameSingle);
	elements.buttons.startBattle.addEventListener('click', startGameBattle);
	
	// Fetch highscore
	var xhr = new XMLHttpRequest();
	xhr.open('GET', '/api/highscores')
	xhr.onload = function() {
	    if (xhr.status === 200) {
	    	
	        var res = JSON.parse(xhr.responseText),
	        	html = '',
	        	
	        	current = 0,
	        	max = 10;
	        	
			if (res.Ath) res.Ath.forEach(function(hs) {
				if (current++ < max) {
					html += "<div class=\"highscore-entry\"><h5 class=\"right no-margin\">" + hs.Score + "</h5><h5 class=\"no-margin\">" + hs.Nickname + "</h5></div>"
				}
			});
			
			elements.containers.hsAth.innerHTML = html;
			
        	html = ''; current = 0; max = 10;
        
        	if (res.Week) res.Week.forEach(function(hs) {
				if (current++ < max) {
					html += "<div class=\"highscore-entry\"><h5 class=\"right no-margin\">" + hs.Score + "</h5><h5 class=\"no-margin\">" + hs.Nickname + "</h5></div>"
				}
			});
			
			elements.containers.hsWeek.innerHTML = html;
			
			
	    }
	};
	xhr.send();

});