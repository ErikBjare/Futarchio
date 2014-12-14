Template.votebuttons.helpers({
    "votedStyle": function(upOrDown) {
        vote = Votes.findOne({post: this._id, createdBy: Meteor.userId(), type: "UpDown"});
        if(vote === undefined) {
            return "";
        }

        if(vote.value === 1) {
            return upOrDown === "up" ? "color: #3F3;" : "";
        } else if (vote.value === -1) {
            return upOrDown === "down" ? "color: #F33;" : "";
        } else {
            console.error("Invalid vote.value: " + vote.value);
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

        var vote = new Vote({
            type: "UpDown",
            value: event.target.id === "up" ? 1 : -1,
            post: template.data._id
        });
        Votes.insert(vote);
    }
});

