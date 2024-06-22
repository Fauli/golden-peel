document.addEventListener('DOMContentLoaded', function() {
    const gallery = document.querySelector('.gallery');

    for (let i = 1; i <= 24; i++) {
        const img = document.createElement('img');
        img.src = `Images/Banana-Candles-${i}.webp`;
        img.alt = `Banana Candles ${i}`;
        gallery.appendChild(img);

        img.addEventListener('click', function() {
            openModal(img.src, img.alt);
        });
    }
});

function openModal(src, alt) {
    const modal = document.createElement('div');
    modal.className = 'modal';
    modal.innerHTML = `
        <span class="close">&times;</span>
        <img class="modal-content" src="${src}">
        <div class="modal-caption">${alt}</div>
    `;

    document.body.appendChild(modal);
    modal.style.display = "block";

    modal.querySelector('.close').addEventListener('click', function() {
        modal.style.display = "none";
        modal.remove();
    });
}
