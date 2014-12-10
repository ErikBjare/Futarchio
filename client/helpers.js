Template.registerHelper("session", function(str) {
    return Session.get(str);
});

Template.registerHelper("usernameOf", function(userid) {
    user = Meteor.users.findOne({"_id": userid});
    if(user === undefined) {
        console.error("User " + userid + " not found");
        return "[ERROR: USER NOT FOUND]";
    }
    if(user.username === undefined) {
        console.error("User " + userid + " has no username");
        return "?";
    }
    return user.username;
});

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

Template.registerHelper("points", function(id) {
    console.log("Fetching points for: " + id);
    upVotes = Votes.find({"post": id, value: 1}).count();
    downVotes = Votes.find({"post": id, value: -1}).count();
    points = upVotes - downVotes;
    return points.toString();
});

