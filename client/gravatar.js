gravatar_hash = function(email) {
    if(!_.contains(email, "@")) {
        console.error("Error: got empty email");
    }
    var hash = CryptoJS.MD5(email).toString();
    return hash;
};
