let count = 0;

function incrementCounter() {
    const banana = document.getElementById('banana');
    const counter = document.getElementById('counter');
    count++;
    counter.innerText = count;

    // Color transition every 100 clicks
    if (count % 100 === 0) {
        let yellowIntensity = Math.min(count / 10000, 1);
        let greenValue = Math.max(255 - yellowIntensity * 255, 50); // Ensuring some green remains
        banana.style.color = `rgb(255, ${greenValue}, 0)`;
    }

    // Wiggle effect after 10,000 clicks
    if (count >= 10000 && count % 100 === 0) {
        banana.style.animation = 'wiggle 0.5s';
        setTimeout(() => {
            banana.style.animation = '';
        }, 500);
    }
}
