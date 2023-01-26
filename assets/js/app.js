(function (w, d){
  'use strict';

  /* Takes a number of seconds and formats into a human readable form like,
   * 90s -> "01:30" */
  function timefmt(seconds) {
    var mins = Math.floor(seconds / 60);

    var minutes = mins.toLocaleString('en-US', {
      minimumIntegerDigits: 2,
      useGrouping: false
    });

    var secs = Math.floor(seconds - (minutes * 60));

    var seconds = secs.toLocaleString('en-US', {
      minimumIntegerDigits: 2,
      useGrouping: false
    });

    return `${minutes}:${seconds}`;
  }

  /* Player is the music player component. The only argument is the root element
   * to mount the player. */
  var Player = function (player) {
    var that = this;

    this.root = player;

    var defaultTrackID = this.root.querySelector('ol > li').dataset.id;

    this.state = {
      mode: 'paused',
      track: 0,
      album: this.root.querySelector('img'),
      tracks: this.root.querySelectorAll('ol > li'),
      playing: this.root.querySelector('p.now_playing'),
      shuffle: this.root.querySelector('input[name="shuffle"]'),
      position: this.root.querySelector('.current_pos'),
      duration: this.root.querySelector('.duration'),
      progress: this.root.querySelector('input[name="progress"]'),
      audio: new Audio(`/track/${defaultTrackID}`),
    };

    this.callbacks = {
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
          that.actions.next();
        },
      ],
      'track_loaded': [
        function (track) {
          that.state.position.innerText = timefmt(that.state.audio.currentTime);
          that.state.duration.innerText = timefmt(that.state.audio.duration);

          that.state.progress.max = that.state.audio.duration;
          that.state.progress.value = that.state.audio.currentTime;
        },
      ],
      'time_update': [
        function (track) {
          that.state.position.innerText = timefmt(that.state.audio.currentTime);
          that.state.duration.innerText = timefmt(that.state.audio.duration);

          that.state.progress.max = that.state.audio.duration;
          that.state.progress.value = that.state.audio.currentTime;
        },
      ],
    };

    this.buttons = {
      play: this.root.querySelector('.player > button.play'),
      next: this.root.querySelector('.player > button.next'),
      prev: this.root.querySelector('.player > button.prev'),
      skip_back: this.root.querySelector('.player > button.skip_back'),
      skip_forward: this.root.querySelector('.player > button.skip_forward'),
    };

    this.actions = {
      dispatch: function (eventName, track) {
        if (that.callbacks[eventName]) {
          that.callbacks[eventName].forEach(function (cb) {
            cb(track);
          });
        }
      },

      skip_back: function () {
        var t = that.state.audio.currentTime - 15;

        if (t < 0) {
          t = 0;
        }

        that.state.position.innerText = t;
        that.state.audio.currentTime = t;
      },

      skip_forward: function () {
        var t = that.state.audio.currentTime + 15;

        if (t > that.state.audio.duration) {
          t = that.state.audio.duration;
        }

        that.state.position.innerText = t;
        that.state.audio.currentTime = t;
      },

      prev: function () {
        that.state.track -= 1;

        if (that.state.track < 0) {
          that.state.track = that.state.tracks.length - 1;
        }

        var track = that.state.tracks[that.state.track];

        that.state.playing.innerText = `${track.dataset.artist} - ${track.dataset.title}`;
        that.state.mode = 'paused';
        that.state.audio.pause();
        that.state.audio.src = `/track/${track.dataset.id}`;
        that.state.album.src = `/thumbnail/${track.dataset.id}`;
        that.state.mode = 'playing';
        that.state.audio.play();
        that.buttons.play.innerHTML = 'Pause';

        that.actions.dispatch('track_changed', track);
      },
      next: function () {
        that.state.track += 1;

        if (that.state.track >= that.state.tracks.length) {
          that.state.track = 0;
        }

        if (that.state.shuffle.checked) {
          that.state.track = Math.floor(Math.random() * that.state.tracks.length);
        }

        var track = that.state.tracks[that.state.track];

        that.state.playing.innerText = `${track.dataset.artist} - ${track.dataset.title}`;
        that.state.mode = 'paused';
        that.state.audio.pause();
        that.state.audio.src = `/track/${track.dataset.id}`;
        that.state.album.src = `/thumbnail/${track.dataset.id}`;
        that.state.mode = 'playing';
        that.state.audio.play();
        that.buttons.play.innerHTML = 'Pause';

        that.actions.dispatch('track_changed', track);
      },
      play: function () {
        that.state.mode = 'playing';
        that.state.audio.play();
        that.buttons.play.innerHTML = 'Pause';
      },
      pause: function () {
        that.state.mode = 'paused';
        that.state.audio.pause();
        that.buttons.play.innerHTML = 'Play';
      },
    };

    /* add click handler for each track in the play list */
    this.state.tracks.forEach(function (track) {
      track.addEventListener('click', function () {
        that.state.track = parseInt(this.dataset.index);
        that.state.mode = 'paused';
        that.state.audio.pause();
        that.state.playing.innerText = `${this.dataset.artist} - ${this.dataset.title}`;
        that.state.audio.src = `/track/${this.dataset.id}`;
        that.state.album.src = `/thumbnail/${this.dataset.id}`;
        that.state.mode = 'playing';
        that.state.audio.play();
        that.buttons.play.innerHTML = 'Pause';

        that.actions.dispatch('track_changed', this);
      });
    });

    /* add click handlers for each player button */
    this.buttons.play.addEventListener('click', function () {
      if (that.state.mode === 'paused') {
        that.actions.play();
      } else {
        that.actions.pause();
      }
    });

    that.buttons.skip_back.addEventListener('click', that.actions.skip_back);
    that.buttons.prev.addEventListener('click', that.actions.prev);
    that.buttons.next.addEventListener('click', that.actions.next);
    that.buttons.skip_forward.addEventListener('click', that.actions.skip_forward);

    /* hook into audio events */
    that.state.audio.onloadedmetadata = function () {
      that.actions.dispatch('track_loaded', that.state.tracks[that.state.track]);
    };

    that.state.audio.ontimeupdate = function () {
      that.actions.dispatch('time_update', that.state.tracks[that.state.track]);
    };

    that.state.audio.onended = function () {
      that.actions.dispatch('track_ended', that.state.tracks[that.state.track]);
    };

    that.state.progress.onchange = function () {
      that.state.position.innerText = this.value;
      that.state.audio.currentTime = this.value;
    };
  }

  var player = new Player(d.querySelector('#player'));
})(window, document)
