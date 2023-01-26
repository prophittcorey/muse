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

    callbacks: {
      'track_changed': [
        function (track) {
          console.log('Track chanegd to ', track.dataset.title);
        },
        function (track) {
          var title = d.querySelector('title');
          title.innerText = `Muse - ${track.dataset.title}`;
        },
      ],
      'track_ended': [
        function (track) {
          player.actions.next();
        },
      ],
      'track_loaded': [
        function (track) {
          var pos = d.querySelector('.current_pos');
          var dur = d.querySelector('.duration');

          pos.innerText = player.state.audio.currentTime;
          dur.innerText = player.state.audio.duration;
        },
      ],
      'time_update': [
        function (track) {
          var pos = d.querySelector('.current_pos');
          var dur = d.querySelector('.duration');

          pos.innerText = player.state.audio.currentTime;
          dur.innerText = player.state.audio.duration;
        },
      ],
    },

    buttons: {
      play: d.querySelector('.player > button.play'),
      next: d.querySelector('.player > button.next'),
      prev: d.querySelector('.player > button.prev'),
    },

    actions: {
      dispatch: function (eventName, track) {
        if (player.callbacks[eventName]) {
          player.callbacks[eventName].forEach(function (cb) {
            cb(track);
          });
        }
      },

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

        player.actions.dispatch('track_changed', track);
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

        player.actions.dispatch('track_changed', track);
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

      player.actions.dispatch('track_changed', this);
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
  player.state.audio.onloadedmetadata = function () {
    player.actions.dispatch('track_loaded', player.state.tracks[player.state.track]);
  };

  player.state.audio.ontimeupdate = function () {
    player.actions.dispatch('time_update', player.state.tracks[player.state.track]);
  };

  player.state.audio.onended = function () {
    player.actions.dispatch('track_ended', player.state.tracks[player.state.track]);
  };
})(window, document)
