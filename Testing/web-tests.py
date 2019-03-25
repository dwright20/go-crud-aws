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
    def signIn(self):
        driver = self.driver
        driver.find_element_by_name("user_name").send_keys(USER)  # send username
        driver.find_element_by_name("user_pass").send_keys(PASS)  # send password
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
    def test_01_signin(self):
        driver = self.driver
        driver.find_element_by_name("user_name").send_keys(USER)
        driver.find_element_by_name("user_pass").send_keys(PASS)
        driver.find_element_by_css_selector('button[type="submit"]').click()
        res = WebDriverWait(driver, 5).until(EC.title_contains('Select'))
        worked = res
        self.assertTrue(worked)  # determine if sign in is successful and taking to correct page

    # test to ensure a cookie is being properly set from web server
    def test_02_cookie(self):
        self.signIn()
        driver = self.driver
        cookie_list = driver.get_cookies()
        worked = True if cookie_list[0].get(u'value') == u'testing' else False
        self.assertTrue(worked)  # determine if value is properly being set to user accessing site

    # test to ensure a cookie is required for all pages past sign-in/creation
    def test_03_req_cookie(self):
        driver = self.driver
        driver.get(WEB_SERVER + "gameSelect")
        worked = True if driver.current_url == WEB_SERVER else False
        self.assertTrue(worked)  # determine if web server is properly redirecting non signed-in users to sign-in page

    # test submitting apex data with test data
    def test_apex_submit(self):
        self.signIn()  # sign into web server
        driver = self.driver
        driver.find_element_by_css_selector('a[href="apexSelect"]').click()  # select Apex game option
        driver.find_element_by_partial_link_text('Submit').click()  # select submission page
        # below finds fields of submission form and sends mock data
        driver.find_element_by_name("result").send_keys('loss')
        driver.find_element_by_name("legend").send_keys("testlegend")
        driver.find_element_by_name("kills").send_keys('0')
        driver.find_element_by_name("placement").send_keys('1')
        driver.find_element_by_name("damage").send_keys('1000')
        driver.find_element_by_name("time").send_keys('15')
        driver.find_element_by_name("teammates").send_keys('testteam')
        driver.find_element_by_css_selector('button[type="submit"]').click()  # submit mock data
        res = WebDriverWait(driver, 5).until(EC.title_contains('Select'))
        worked = res
        self.assertTrue(worked)  # determine if you are taken to correct page after submission

    # test submitting fort data with test data
    def test_fort_submit(self):
        self.signIn()  # sign into web server
        driver = self.driver
        driver.find_element_by_css_selector('a[href="fortniteSelect"]').click()  # select Apex game option
        driver.find_element_by_partial_link_text('Submit').click()  # select submission page
        # below finds fields of submission form and sends mock data
        driver.find_element_by_name("result").send_keys('loss')
        driver.find_element_by_name("kills").send_keys('0')
        driver.find_element_by_name("placement").send_keys('1')
        driver.find_element_by_name("mode").send_keys('duo')
        driver.find_element_by_name("teammates").send_keys('random')
        driver.find_element_by_css_selector('button[type="submit"]').click()  # submit mock data
        res = WebDriverWait(driver, 5).until(EC.title_contains('Select'))
        worked = res
        self.assertTrue(worked)  # determine if you are taken to correct page after submission

    # test submitting hots data with test data
    def test_hots_submit(self):
        self.signIn()  # sign into web server
        driver = self.driver
        driver.find_element_by_css_selector('a[href="hotsSelect"]').click()  # select Apex game option
        driver.find_element_by_partial_link_text('Submit').click()  # select submission page
        # below finds fields of submission form and sends mock data
        driver.find_element_by_name("result").send_keys('loss')
        driver.find_element_by_name("hero").send_keys("testhero")
        driver.find_element_by_name("kills").send_keys('0')
        driver.find_element_by_name("deaths").send_keys('1')
        driver.find_element_by_name("assists").send_keys('10')
        driver.find_element_by_name("time").send_keys('15')
        driver.find_element_by_name("map").send_keys('testmap')
        driver.find_element_by_css_selector('button[type="submit"]').click()  # submit mock data
        res = WebDriverWait(driver, 5).until(EC.title_contains('Select'))
        worked = res
        self.assertTrue(worked)  # determine if you are taken to correct page after submission

    # test viewing apex data
    def test_apex_view(self):
        self.signIn()  # sign into web server
        driver = self.driver
        driver.find_element_by_css_selector('a[href="apexSelect"]').click()
        driver.find_element_by_partial_link_text('View').click()  # select view page
        table = driver.find_element_by_css_selector('tbody').get_attribute('innerHTML')  # read table body html
        worked = True if USER and 'apex' in table else False
        self.assertTrue(worked)  # determine if table is presenting data for account

    # test viewing fort data
    def test_fort_view(self):
        self.signIn()  # sign into web server
        driver = self.driver
        driver.find_element_by_css_selector('a[href="fortniteSelect"]').click()
        driver.find_element_by_partial_link_text('View').click()  # select view page
        table = driver.find_element_by_css_selector('tbody').get_attribute('innerHTML')  # read table body html
        worked = True if USER and 'fort' in table else False
        self.assertTrue(worked)  # determine if table is presenting data for account


    # test viewing hots data
    def test_hots_view(self):
        self.signIn()  # sign into web server
        driver = self.driver
        driver.find_element_by_css_selector('a[href="hotsSelect"]').click()
        driver.find_element_by_partial_link_text('View').click()  # select view page
        table = driver.find_element_by_css_selector('tbody').get_attribute('innerHTML')  # read table body html
        worked = True if USER and 'hots' in table else False
        self.assertTrue(worked)  # determine if table is presenting data for account

    # close chrome session after each test
    def tearDown(self):
        self.driver.close()

if __name__== '__main__':
    unittest.main()
