document.addEventListener('DOMContentLoaded', () => {
    const btnPlay = document.getElementById('play');
    const btnPrev = document.getElementById('previous');
    const btnNext = document.getElementById('next');
    const links = document.querySelectorAll('a');
    let playing = false;

    function handler(e) {
        e.preventDefault();

        if (playing) {
            btnPlay.innerHTML = '►';
            playing.pause();
            playing.currentTime = 0
        }
        let file = e.currentTarget.href || document.querySelectorAll('li > a')[0].href;
        playing = new Audio(file);
        playing.addEventListener('timeupdate', e => {
            document.getElementById('time').style.width = `${playing.currentTime / playing.duration * 100}%`
        });
        playing.play();
        btnPlay.innerHTML = '❚❚';
    }

    function nextHandler(e) {
        e.preventDefault();
        skip(1);
    }

    function prevHandler(e) {
        e.preventDefault();
        skip(-1);
    }

    function skip(dir) {
        if (!playing) {
            return;
        }
        let current = document.querySelector(`[href="${playing.src}"]`);
        let next;
        if (dir > 0) {
            next = current.parentNode.nextSibling.firstChild;
        } else if (dir < 0) {
            next = current.parentNode.previousSibling.firstChild;
        }
        if (!next) {
            return;
        }
        playing.pause();
        playing.currentTime = 0;
        playing = new Audio(next.href);
        playing.play();
        btnPlay.innerHTML = '❚❚';
    }

    function playHandler(e) {
        e.preventDefault();

        if (!playing) {
            return;
        }

        if (!playing.paused) {
            playing.pause();
            btnPlay.innerHTML = '►';
        } else {
            playing.play();
            btnPlay.innerHTML = '❚❚';
        }
    }

    for (let i = 0; i < links.length; i++) {
        let link = links[i];
        link.addEventListener('click', handler);
    }
    btnPlay.addEventListener('click', playHandler);
    btnNext.addEventListener('click', nextHandler);
    btnPrev.addEventListener('click', prevHandler);
})