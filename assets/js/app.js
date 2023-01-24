(function (w, d){
  'use strict';

  var player = {
    state: {
      album: d.querySelector('main > img'),
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
        player.buttons.play.innerHTML = 'Pause';
      },
      pause: function () {
        player.state.mode = 'paused';
        player.state.audio.pause();
        player.buttons.play.innerHTML = 'Play';
      },
    },
  };

  /* add click handler for each track in the play list */
  player.state.tracks.forEach(function (track) {
    track.addEventListener('click', function () {
      player.state.mode = 'paused';
      player.state.audio.pause();
      player.state.audio.src = `/track/${this.dataset.id}`;
      player.state.album.src = `/thumbnail/${this.dataset.id}`;
      player.state.mode = 'playing';
      player.state.audio.play();
      player.buttons.play.innerHTML = 'Pause';
    });
  });

  /* add click handlers for each player button */
  player.buttons.play.addEventListener('click', function () {
    if (player.state.mode === 'paused') {
      player.actions.play();
    } else {
      player.actions.pause();
    }
  });
})(window, document)
