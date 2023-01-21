(function (w, d){
  'use strict';

  // TODO: Try next/previous/play/pause actual audio.

  var player = {
    state: {
      track: 0,
      tracks: d.querySelectorAll('main > ol > li'),
      audio: new Audio(`/track/${d.querySelector('main > ol > li').dataset.id}`),
      mode: 'paused',
    },

    buttons: {
      play: d.querySelector('.player > button'),
    },

    actions: {
      play: function () {
        player.state.mode = 'playing';
        player.state.audio.play();
      },
      pause: function () {
        player.state.mode = 'paused';
        player.state.audio.pause();
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
