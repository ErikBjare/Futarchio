capitalize = function(string) {
    return string[0].toUpperCase() + string.slice(1);
};

// TODO: Find a better way, make it work server-side if needed (is it?)
// From: https://stackoverflow.com/questions/10984030/get-meteor-collection-by-name
getCollection = function(string) {
    for (var globalObject in window) {
        if (window[globalObject] instanceof Meteor.Collection) {
            if (globalObject === string) {
                return (window[globalObject]);
            }
        }
    }
    return undefined; // if none of the collections match
};
