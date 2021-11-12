$(function () {
    appHeader().stickyOnScroll()
})


const appHeader = function () {
    const header = {
        fixed: false,
        scrollTopLimit: 360,
        $element: $("#header"),
    }

    header.stickyOnScroll = function () {
        let self = this
        $(window).scroll(function () {
            let scroll = $(window).scrollTop();
            if (scroll > self.scrollTopLimit && !self.fixed) {
                self.$element.addClass("js-header-fixed")
                self.fixed = true
            } else if (scroll < self.scrollTopLimit && self.fixed) {
                self.$element.removeClass("js-header-fixed")
                self.fixed = false
            }
        })
    }

    return header
}