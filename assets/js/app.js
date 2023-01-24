(function (w, d){
  'use strict';

  var player = {
    state: {
      album: d.querySelector('main > img'),
      track: 0,
      tracks: d.querySelectorAll('main > ol > li'),
      playing: d.querySelector('p.now_playing'),
      shuffle: d.querySelector('input[name="shuffle"]'),
      audio: new Audio(`/track/${d.querySelector('main > ol > li').dataset.id}`),
      mode: 'paused',
    },

    buttons: {
      play: d.querySelector('.player > button.play'),
      next: d.querySelector('.player > button.next'),
      prev: d.querySelector('.player > button.prev'),
    },

    actions: {
      prev: function () {
        player.state.track -= 1;

        if (player.state.track < 0) {
          player.state.track = player.state.tracks.length - 1;
        }

        var track = player.state.tracks[player.state.track];

        player.state.playing.innerText = `${track.dataset.artist} - ${track.dataset.title}`;
        player.state.mode = 'paused';
        player.state.audio.pause();
        player.state.audio.src = `/track/${track.dataset.id}`;
        player.state.album.src = `/thumbnail/${track.dataset.id}`;
        player.state.mode = 'playing';
        player.state.audio.play();
        player.buttons.play.innerHTML = 'Pause';
      },
      next: function () {
        player.state.track += 1;

        if (player.state.track >= player.state.tracks.length) {
          player.state.track = 0;
        }

        if (player.state.shuffle.checked) {
          player.state.track = Math.floor(Math.random() * player.state.tracks.length);
        }

        var track = player.state.tracks[player.state.track];

        player.state.playing.innerText = `${track.dataset.artist} - ${track.dataset.title}`;
        player.state.mode = 'paused';
        player.state.audio.pause();
        player.state.audio.src = `/track/${track.dataset.id}`;
        player.state.album.src = `/thumbnail/${track.dataset.id}`;
        player.state.mode = 'playing';
        player.state.audio.play();
        player.buttons.play.innerHTML = 'Pause';
      },
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
      player.state.track = parseInt(this.dataset.index);
      player.state.mode = 'paused';
      player.state.audio.pause();
      player.state.playing.innerText = `${this.dataset.artist} - ${this.dataset.title}`;
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

  player.buttons.prev.addEventListener('click', player.actions.prev);
  player.buttons.next.addEventListener('click', player.actions.next);

  /* hook into audio events */
  player.state.audio.onended = function () {
    player.actions.next();
  };
})(window, document)
