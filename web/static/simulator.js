const runSimulator = () => {
  let isActive = false;

  if (!isActive) return;

  fetch('static/seed.json')
    .then((r) => r.json())
    .then((data) => {
      const positives = data.phrases.positive;
      const negatives = data.phrases.negative;

      const size =
        positives.length > negatives.length
          ? positives.length
          : negatives.length;

      let input = document.getElementById('text');
      let btn = document.getElementById('send');
      for (let i = 0; i < size; i++) {
        setTimeout(() => {
          if (positives[i]) {
            input.value = positives[i];
            btn.click();
          }
        }, 1000 * i);

        setTimeout(() => {
          if (negatives[i]) {
            input.value = negatives[i];
            btn.click();
          }
        }, 2000 * i);
      }
    });
};
