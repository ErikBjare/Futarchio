Template.votebuttons.helpers({
    "votedStyle": function(upOrDown) {
        if(!Meteor.userId()) {
            return "";
        }

        var rating = this.ratings[Meteor.userId()];
        if(!rating) {
            return "";
        }

        if(rating.score === 1) {
            return upOrDown === "up" ? "color: #3F3;" : "";
        } else if (rating.score === -1) {
            return upOrDown === "down" ? "color: #F33;" : "";
        } else {
            console.error("Invalid score: " + rating.score);
            return "";
        }
    }
});

Template.votebuttons.events({
    "click": function(event, template) {
        if(!_.contains(["up", "down"], event.target.id)) {
            console.warn("invalid votebutton id: " + event.target.id);
            return;
        }

        var collectionName = template.data.cardType[0].toUpperCase() + template.data.cardType.slice(1) + "s";
        var collection = getCollection(collectionName);

        var score = event.target.id === "up" ? 1 : -1;

        var insert;
        if(template.data.ratings[Meteor.userId()] !== undefined &&
            template.data.ratings[Meteor.userId()].score === score) {
            insert = {$unset: {}};
            insert.$unset["ratings."+Meteor.userId()] = "";
        } else {
            insert = {$set: {}};
            insert.$set["ratings."+Meteor.userId()] = new Rating(score);
        }

        var n = collection.update(template.data._id, insert);
        if(n != 1) {
            console.warn("Updated " + n + " " + collectionName);
        }
    }
});

