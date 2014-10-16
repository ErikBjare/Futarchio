app.filter('capitalize', function() {
    return function(input) {
        if (input && input[0]) {
            output = input[0].toUpperCase() + input.substr(1,input.length);
            return output;
        }
        return "";
    };
});
