$(function(){
    $('.editormd-toc-menu').mouseover(function(){
        $(this).children('.editormd-markdown-toc').children('.markdown-toc-list').show();
    }).mouseleave(function(){
        $(this).children('.editormd-markdown-toc').children('.markdown-toc-list').hide();
    });

    var timeSince;

    timeSince = function(date) {
        var interval, seconds;
        seconds = Math.floor((new Date() - date) / 1000);
        interval = Math.floor(seconds / 31536000);
        if (interval > 1) {
            return interval + " years ago";
        }
        interval = Math.floor(seconds / 2592000);
        if (interval > 1) {
            return interval + " months ago";
        }
        interval = Math.floor(seconds / 86400);
        if (interval > 1) {
            return interval + " days ago";
        }
        interval = Math.floor(seconds / 3600);
        if (interval > 1) {
            return interval + " hours ago";
        }
        interval = Math.floor(seconds / 60);
        if (interval > 1) {
            return interval + " mins ago";
        }
        return Math.floor(seconds) + " seconds ago";
    };

    $('.search').bind($.clickOrTouch(),function(){
        $.dialog({
            type: 'table',
            from : '#search-form',
            padding: 30
        });
    });

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
            $('#go-top').show();
        else
            $('#go-top').hide();
    });

    $('#go-top').click(function(){
        $('html,body').animate({
            scrollTop: '0px'
        }, 800);
    });

    $('img').each(function(idx, item) {
      var $item, imageAlt;
      $item = $(item);
      if ($item.attr('data-src')) {
        $item.wrap('<a href="' + $item.attr('data-src') + '" target="_blank"></a>');
        imageAlt = $item.prop('alt');
        if ($.trim(imageAlt)) {
          return $item.parent('a').after('<div class="image-alt">' + imageAlt + '</div>');
        }
      }
    });

    if ($('img').unveil) {
      return $('img').unveil(200, function() {
        return $(this).load(function() {
          return this.style.opacity = 1;
        });
      });
    }
});