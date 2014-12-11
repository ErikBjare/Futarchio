Template.card.created = function() {
    this.data.showDetails = new ReactiveVar();
    this.data.showDetails.set(false);
};

Template.card.helpers({
    error: "",
    showDetails: function() {
        return this.showDetails.get();
    },
    detailsTemplate: function() {
        return Template[Template.parentData().type + "Details"];
    }
});

Template.card.events({
    "click #showDetails": function(e, template) {
        template.data.showDetails.set(!(template.data.showDetails.get()));
    }
});
