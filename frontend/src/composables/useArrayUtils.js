
export function useArrayUtils() {
    /**
     * Shuffles array in place. Fisher–Yates shuffle algorithm ES6 version.
     * @param {Array} a items An array containing the items.
     */
    function shuffle(a) {
        if (!a) return a

        for (let i = a.length - 1; i > 0; i--) {
            const j = Math.floor(Math.random() * (i + 1));
            [a[i], a[j]] = [a[j], a[i]];
        }
        return a;
    }

    return { shuffle}
}
