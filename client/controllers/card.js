Template.card.created = function() {
    this.data.showDetails = new ReactiveVar();
    this.data.showDetails.set(false);

    // TODO: Ugly hack, make it beautiful and safe
    var cardType = this.view.parentView.parentView.name.split(".")[1];
    if(Template[cardType + "Details"] === undefined) {
        console.error("Card details template '" + cardType + "Details', doesn't exist");
    }
    this.data.cardType = cardType;
};

Template.card.helpers({
    error: "",
    showDetails: function() {
        return this.showDetails.get();
    },
    detailsTemplate: function() {
        return Template[this.cardType + "Details"];
    }
});

Template.card.events({
    "click #showDetails": function(e, template) {
        template.data.showDetails.set(!(template.data.showDetails.get()));
    }
});
