app.filter('capitalize', function() {
    return function(input) {
        if (input && input[0]) {
            output = input[0].toUpperCase() + input.substr(1,input.length);
            return output;
        }
        return "";
    };
});

app.filter('fromNow', function() {
    return function(datetimeStr) {
        return moment(datetimeStr).fromNow();
    };
});
