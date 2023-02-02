(function (w, d){
  'use strict';

  /* Takes a number of seconds and formats into a human readable form like,
   * 90s -> "01:30" */
  function timefmt(seconds) {
    var options = {
      minimumIntegerDigits: 2,
      useGrouping: false
    };

    var mins = Math.floor(seconds / 60);
    var minutes = mins.toLocaleString('en-US', options);

    var secs = Math.floor(seconds - (minutes * 60));
    var seconds = secs.toLocaleString('en-US', options);

    return `${minutes}:${seconds}`;
  }

  /* Player is the music player component. The only argument is the root element
   * to mount the player. */
  var Player = function (player) {
    this.root = player;

    var _player = this;

    this.state = {
      mode: 'paused',
      track: 0,
      album: this.root.querySelector('img'),
      tracks: this.root.querySelectorAll('ol > li'),
      playing: this.root.querySelector('.now_playing'),
      artist: this.root.querySelector('.current_artist'),
      shuffle: this.root.querySelector('input[name="shuffle"]'),
      repeat: this.root.querySelector('input[name="repeat"]'),
      position: this.root.querySelector('.current_pos'),
      duration: this.root.querySelector('.duration'),
      progress: this.root.querySelector('input[name="progress"]'),
      title: d.querySelector('title'),
      audio: new Audio(`/track/${this.root.querySelector('ol > li').dataset.id}`),
    };

    this.callbacks = {
      'track_changed': [
        function (track) { _player.state.title.innerText = `Muse - ${track.dataset.title}`; },
      ],
      'track_ended': [
        function (track) { _player.actions.move(1); },
      ],
      'track_loaded': [
        function (track) {
          if (isNaN(_player.state.audio.duration) || isNaN(_player.state.audio.currentTime)) {
            return;
          }

          _player.state.position.innerText = timefmt(_player.state.audio.currentTime);
          _player.state.duration.innerText = timefmt(_player.state.audio.duration);

          _player.state.progress.max = _player.state.audio.duration;
          _player.state.progress.value = _player.state.audio.currentTime;
        },
        function (track) {
          d.querySelectorAll('.active').forEach(function (el) {
            el.classList.remove('active');
          });

          if (!track.classList.contains('active')) {
            track.classList.add('active');
          }
        },
      ],
      'time_update': [
        function (track) {
          if (isNaN(_player.state.audio.duration) || isNaN(_player.state.audio.currentTime)) {
            return;
          }

          _player.state.position.innerText = timefmt(_player.state.audio.currentTime);
          _player.state.duration.innerText = timefmt(_player.state.audio.duration);

          _player.state.progress.max = _player.state.audio.duration;
          _player.state.progress.value = _player.state.audio.currentTime;
        },
      ],
    };

    this.buttons = {
      play: this.root.querySelector('.player .play'),
      next: this.root.querySelector('.player .next'),
      prev: this.root.querySelector('.player .prev'),
    };

    this.actions = {
      dispatch: function (eventName, track) {
        if (_player.callbacks[eventName]) {
          _player.callbacks[eventName].forEach(function (cb) {
            cb(track);
          });
        }
      },

      move: function (direction) {
        if (_player.state.repeat.checked) {
          direction = 0
        }

        _player.state.track += direction;

        if (_player.state.track < 0) {
          _player.state.track = _player.state.tracks.length - 1;
        }

        if (_player.state.track >= _player.state.tracks.length) {
          _player.state.track = 0;
        }

        if (_player.state.shuffle.checked) {
          _player.state.track = Math.floor(Math.random() * _player.state.tracks.length);
        }

        var track = _player.state.tracks[_player.state.track];

        track.scrollIntoView({ behavior: 'auto', block: 'center', inline: 'center' });

        _player.state.playing.innerText = track.dataset.title;
        _player.state.artist.innerText = track.dataset.artist;
        _player.state.mode = 'paused';
        _player.state.audio.pause();
        _player.state.audio.src = `/track/${track.dataset.id}`;
        _player.state.album.src = `/thumbnail/${track.dataset.id}`;
        _player.state.mode = 'playing';
        _player.state.audio.play();
        _player.buttons.play.querySelector('img').src = '/assets/images/pause.svg';

        _player.actions.dispatch('track_changed', track);
      },

      toggle: function () {
        _player.state.mode === 'paused' ? _player.actions.play() : _player.actions.pause();
      },

      play: function () {
        _player.state.mode = 'playing';
        _player.state.audio.play();
        _player.buttons.play.querySelector('img').src = '/assets/images/pause.svg';
      },

      pause: function () {
        _player.state.mode = 'paused';
        _player.state.audio.pause();
        _player.buttons.play.querySelector('img').src = '/assets/images/play.svg';
      },
    };

    /* add click handler for each track in the play list */
    this.state.tracks.forEach(function (track) {
      track.addEventListener('click', function () {
        _player.state.track = parseInt(this.dataset.index);
        _player.state.mode = 'paused';
        _player.state.audio.pause();
        _player.state.playing.innerText = this.dataset.title;
        _player.state.artist.innerText = this.dataset.artist;
        _player.state.audio.src = `/track/${this.dataset.id}`;
        _player.state.album.src = `/thumbnail/${this.dataset.id}`;
        _player.state.mode = 'playing';
        _player.state.audio.play();
        _player.buttons.play.querySelector('img').src = '/assets/images/pause.svg';

        _player.actions.dispatch('track_changed', this);
      });
    });

    /* add click handlers for each player button */
    this.buttons.play.addEventListener('click', function () { _player.actions.toggle(); });

    _player.buttons.prev.addEventListener('click', function () { _player.actions.move(-1);  });
    _player.buttons.next.addEventListener('click', function () { _player.actions.move(1);   });

    /* hook into audio events */
    _player.state.audio.onloadedmetadata = function () {
      _player.actions.dispatch('track_loaded', _player.state.tracks[_player.state.track]);
    };

    _player.state.audio.ontimeupdate = function () {
      _player.actions.dispatch('time_update', _player.state.tracks[_player.state.track]);
    };

    _player.state.audio.onended = function () {
      _player.actions.dispatch('track_ended', _player.state.tracks[_player.state.track]);
    };

    _player.state.progress.oninput = function () {
      _player.state.position.innerText = timefmt(this.value);
      _player.state.audio.currentTime = this.value;
    };

    d.onkeydown = function (e) {
      e = e || w.event;

      switch (e.key) {
        case 'Enter':
        case 'c':
        case ' ':
          _player.actions.toggle();
          break;

        case 'ArrowLeft':
        case 'h':
        case 'k':
          _player.actions.move(-1);
          break;

        case 'ArrowRight':
        case 'j':
        case 'l':
          _player.actions.move(1);
          break;
      }
    };

    if (navigator && navigator.mediaSession && navigator.mediaSession.setActionHandler) {
      navigator.mediaSession.setActionHandler('previoustrack', function () {
        _player.actions.move(-1);
      });

      navigator.mediaSession.setActionHandler('nexttrack', function () {
        _player.actions.move(1);
      });

      navigator.mediaSession.setActionHandler('play', function () {
        _player.actions.toggle();
      });

      navigator.mediaSession.setActionHandler('pause', function () {
        _player.actions.toggle();
      });
    }
  };

  new Player(d.querySelector('.player'));
})(window, document)
