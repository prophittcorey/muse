(function (w, d){
  'use strict';

  // TODO: Add an Audio element..
  // TODO: Try next/previous/play/pause actual audio.
  // TODO: Make audio handler ("/audio/:id").

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
