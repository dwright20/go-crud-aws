import unittest
from selenium import webdriver
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC

DRIVER = ''
WEB_SERVER = ''
USER = ''  # testing username
PASS = ''  # testing password


#  tests for web server
class WebTesting (unittest.TestCase):
    worked = None  # bool used for evaluating success of each test

    # method to sign into web server for testing features past sign in
    def Signin(self):
        driver = self.driver
        username = driver.find_element_by_name("user_name")  # find username field
        password = driver.find_element_by_name("user_pass")  # find password field
        username.send_keys(USER)  # send username
        password.send_keys(PASS)  # send password
        driver.find_element_by_css_selector('button[type="submit"]').click()
        WebDriverWait(driver, 5).until(EC.title_contains('Select'))  # let new webpage load

    # create chrome session and access web server before each test
    def setUp(self):
        # create a new Chrome session
        self.driver = webdriver.Chrome(DRIVER)
        self.driver.implicitly_wait(30)
        # navigate to the web server
        self.driver.get(WEB_SERVER)

    # test sign in functionality with testing credentials
    def test_signin(self):
        driver = self.driver
        username = driver.find_element_by_name("user_name")
        password = driver.find_element_by_name("user_pass")
        username.send_keys(USER)
        password.send_keys(PASS)
        driver.find_element_by_css_selector('button[type="submit"]').click()
        res = WebDriverWait(driver, 5).until(EC.title_contains('Select'))
        worked = res
        self.assertTrue(worked)  # determine if sign in is successful and taking to correct page

    # test submitting apex data with test data
    def test_apex_submit(self):
        self.Signin()  # sign into web server
        driver = self.driver
        driver.find_element_by_css_selector('a[href="apexSelect"]').click()  # select Apex game option
        driver.find_element_by_partial_link_text('Submit').click()  # select submission page
        # below finds fields of submission form and sends mock data
        username = driver.find_element_by_name("user_name")
        result = driver.find_element_by_name("result")
        legend = driver.find_element_by_name("legend")
        kills = driver.find_element_by_name("kills")
        placement = driver.find_element_by_name("placement")
        damage = driver.find_element_by_name("damage")
        time = driver.find_element_by_name("time")
        teammates = driver.find_element_by_name("teammates")
        username.send_keys(USER)
        result.send_keys('loss')
        legend.send_keys("testhero")
        kills.send_keys('0')
        placement.send_keys('1')
        damage.send_keys('1000')
        time.send_keys('15')
        teammates.send_keys('testteam')
        driver.find_element_by_css_selector('button[type="submit"]').click()  # submit mock data
        res = WebDriverWait(driver, 5).until(EC.title_contains('Select'))
        worked = res
        self.assertTrue(worked)  # determine if you are taken to correct page after submission

    # test viewing apex data
    def test_apex_view(self):
        self.Signin()  # sign into web server
        driver = self.driver
        driver.find_element_by_css_selector('a[href="apexSelect"]').click()
        driver.find_element_by_partial_link_text('View').click()  # select view page
        username = driver.find_element_by_name("user_name")
        username.send_keys(USER)
        driver.find_element_by_css_selector('button[type="submit"]').click()
        table = driver.find_element_by_css_selector('tbody')  # find table body
        table = table.get_attribute('innerHTML')  # read table body html
        worked = True if USER in table else False
        self.assertTrue(worked)  # determine if table is presenting data for account

    # close chrome session after each test
    def tearDown(self):
        self.driver.close()

if __name__== '__main__':
    unittest.main()
