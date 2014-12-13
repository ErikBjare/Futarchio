Template.votebuttons.events({
    "click": function(event, template) {
        if(!_.contains(["up", "down"], event.target.id)) {
            console.warn("invalid votebutton id: " + event.target.id);
            return;
        }

        var vote = new Vote({
            type: "UpDown",
            value: event.target.id == "up" ? 1 : -1,
            post: template.data._id
        });
        Votes.insert(vote);
    }
});

