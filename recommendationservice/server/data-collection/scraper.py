import json
from selenium import webdriver
from msedge.selenium_tools import EdgeOptions
import csv
import re

CUISINES = []

class Scraper():
    def __init__(self, file, output):
        self.data = []
        self.filename = file
        self.success_count = 0
        self.count = 0
        self.output_filename = output


    def scrape_website(self,url):
        print(url)
        options = EdgeOptions()
        options.use_chromium = True
        options.add_argument("headless")
        options.add_experimental_option('excludeSwitches', ['enable-logging'])
        driver = webdriver.Chrome(r"C:\Users\benja\Documents\Code\msedgedriver.exe", options=options)
        driver.get(url)

        obj = {}
        # get id=review

        # class=review-title: __ Review: ... -> get name before " Review:" for name
        try:
            title = driver.find_element_by_xpath('//h5[@class="review-title"]')
            title_str = self.extract_prefix(" Review:", title.text)
            if title_str != -1:
                obj["title"] = title_str
            
        
            # class=review-total-box: __/10 -> get before /10, divide by 2 and round down for rating
            rating = driver.find_element_by_xpath('//span[@class="review-total-box"]')
            rating_str = self.extract_prefix("/10", rating.text)
            if rating_str != -1:
                obj["rating"] = int(float(rating_str) // 2)


            # class=post-entry: search for tags wrt to preset array, if within preset output for cuisine
            

            # class=post-entry: search for all instances of $, average out and output for budget
            elements_with_price = driver.find_elements_by_xpath("//*[contains(text(), '$') and not(self::option) and not(self::a)]")
            prices = []
            for element in elements_with_price:
                pattern1 = r'\(\$([0-9]+(?:\.[0-9]+)?)\)'
                matches1 = re.findall(pattern1, element.text)

                # Pattern to match "$__ " at the end of the string
                pattern2 = r'\$([0-9]+(?:\.[0-9]+)?)(?=\s|\+|$)'
                matches2 = re.findall(pattern2, element.text)

                all_matches = matches1 + matches2
                prices += all_matches

            if prices:
                avg = sum(float(price) for price in prices) / len(prices)
                budget = round(avg)
                obj["budget"] = budget


            # class=post-tags: search for all tags, output as array for tags
            tags = driver.find_elements_by_xpath('//div[@class="post-tags"]/a')
            text_list = [tag.text.lower() for tag in tags]
            if text_list:
                obj["tags"] = text_list


            if obj:
                self.data.append(obj)
                self.success_count += 1
            print(self.data)
        except Exception as e:
            print("Unable to locate elements")


    def extract_prefix(self, postfix, string):
        index = string.find(postfix)
        if index != -1:
            return string[:index]
        
        return -1


    def output_file(self):
        field_names = list(self.data[0].keys())

        # append mode: 'a' instead of 'w'
        with open(self.output_filename, 'a', newline='') as file:
            csv_writer = csv.DictWriter(file, fieldnames=field_names)
            csv_writer.writeheader()
            csv_writer.writerows(self.data)


    def run(self):
        with open(self.filename, 'r') as json_file:
            urls = json.load(json_file)

        
        try:
            # NOTE: continue from here
            self.count = 561
            for i in range(561, len(urls)):
                self.count += 1
                self.scrape_website(urls[i])
            # for url in urls:
            #     self.count += 1
            #     self.scrape_website(url)
            # for i in range(2):
            #     self.scrape_website(urls[i])
            
            print(f"Successfully attained {self.success_count} out of {len(urls)} scraped data")
            print("Output to file....")
            
        except KeyboardInterrupt:
            print("\nScript interrupted by user.")
            print(f"Processed {self.count} out of {len(urls)} and successfully extracted {self.success_count}")
            print("Outputting data scraped so far to file....")
        except Exception as e:
            print(f"An error occurred during scraping: {str(e)}. Processed {self.count} out of {len(urls)} URLs")
            print("Outputting data scraped so far to file....")
        
        self.output_file()  
        

    
def main():
    s = Scraper("discovered_urls.json", "output.csv")
    s.run()

if __name__ == "__main__":
    main()