

// maybe add this kind of function in a middleware or in the globalstate
function fetchData(url) {
    return new Promise((resolve, reject) => {
        fetch(url, {
            method: "GET"
        })
            .then(response => {
                return response.json();
            })
            .then(data => {
                console.log(data);
                resolve(data);
            })
            .catch(err => {
                console.log(err);
                reject(err);
            })
    })
}

export default fetchData;