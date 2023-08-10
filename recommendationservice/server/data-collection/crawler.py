import scrapy
import json

class Crawler(scrapy.Spider):
    name = "go_eats_crawler"
    start_urls = ["https://eatbook.sg/category/food-reviews/"]
    allowed_domain = ["https://eatbook.sg/category/food-reviews/", "https://eatbook.sg/"]
    disallowed_domains = ["https://eatbook.sg/wp-content/", "https://eatbook.sg/author/", "https://eatbook.sg/tag/", "http://web.archive.org/", "whatsapp://", "https://www.facebook.com/", "https://twitter.com/"]
    dont_enter = ["https://eatbook.sg/category/"]
    discovered_urls = set()
    max_sites = 750


    def parse(self, response):
        if len(self.discovered_urls) >= self.max_sites:
            self.crawler.engine.close_spider(self, "Reached maximum number of sites to crawl.")
            return

        for link in response.css('a::attr(href)').getall():
            if self.is_within_allowed_domain(link) and not self.is_in_disallowed(link):
                if not self.is_not_enter(link):
                    self.discovered_urls.add(link)
                yield response.follow(link, callback=self.parse)

    def is_within_allowed_domain(self, url):
        for domain in self.allowed_domain:
            if domain in url:
                return True
        return False
    
    def is_not_enter(self, url):
        for domain in self.dont_enter:
            if domain in url:
                return True
        return False
    
    def is_in_disallowed(self, url):
        for domain in self.disallowed_domains:
            if domain in url:
                return True
            
        return False

    def closed(self, reason):
        with open('discovered_urls.json', 'w') as json_file:
            json.dump(list(self.discovered_urls), json_file)