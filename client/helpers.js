Template.registerHelper("capitalize", function(str) {
    if(str !== undefined && str.length > 0) return str[0].toUpperCase() + str.substr(1);
});

Template.registerHelper("fromNow", function(date) {
    if(date === undefined) {
        console.error("date missing");
        return "";
    } else {
        return moment(date).fromNow();
    }
});
