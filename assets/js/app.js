(function (w, d){
  'use strict';

  var player = {
    state: {
      track: 0,
      tracks: d.querySelectorAll('main > ol > li'),
      mode: 'paused',
    },

    buttons: {
      play: d.querySelector('.player > button'),
    },

    actions: {
      play: function () {
        player.state.mode = 'playing';
      },
      pause: function () {
        player.state.mode = 'paused';
      },
    },
  };

  player.buttons.play.addEventListener('click', function () {
    if (player.state.mode === 'paused') {
      player.actions.play();
    } else {
      player.actions.pause();
    }
  });
})(window, document)
