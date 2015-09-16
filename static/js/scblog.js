$(function(){
    var timeSince;

    timeSince = function(date) {
        var interval, seconds;
        seconds = Math.floor((new Date() - date) / 1000);
        interval = Math.floor(seconds / 31536000);
        if (interval > 1) {
            return interval + " 年前";
        }
        interval = Math.floor(seconds / 2592000);
        if (interval > 1) {
            return interval + " 月前";
        }
        interval = Math.floor(seconds / 86400);
        if (interval > 1) {
            return interval + " 天前";
        }
        interval = Math.floor(seconds / 3600);
        if (interval > 1) {
            return interval + " 小时前";
        }
        interval = Math.floor(seconds / 60);
        if (interval > 1) {
            return interval + " 分钟前";
        }
        return Math.floor(seconds) + " 秒前";
    };

    $('.date').each(function(idx, item) {
        var $date, date, timeStr, unixTime;
        $date = $(item);
        timeStr = $date.data('time');
        if (timeStr) {
            unixTime = Number(timeStr) * 1000;
            date = new Date(unixTime);
            return $date.prop('title', date).find('.time').text(timeSince(date));
        }
    });

    $(window).bind('scroll',function(){
        var scrollTop = $(window).scrollTop();
        if(scrollTop > 100)
            $('#gotop').show();
        else
            $('#gotop').hide();
    });

    $('#gotop').click(function(){
        $('html,body').animate({
            scrollTop: '0px'
        }, 800);
    });

    $('#search').click(function(){
        $('.search').toggle(500);
    });
});