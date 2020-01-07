Smart Card Test Tool
===================

Smart Card Test Tool is simple Web Service written in GO for testing PIV Applets in Smart Card.  For more info search for ISO7816-4.

----------

Instructions
-------------

Simply start smartcard.exe and open http://localhost:5580 to open test page.
If you are using PostMan, you will find API collection definition with all available REST URL's

> **Note:**

> - Data parameter must be set in HEX through web page.
> - Data parameter must be set in Base64 encoded HEX array when API is used directly.
> - Response Data property is HEX String.


#### <i class="icon-file"></i> SmartCard Info

SmartCard with PIV applets store client SSL certificate which can be used to authenticate user to web site.

Type of smart card devices are

  - Mobile SIM
  - Bank / Credit Card
  - Security Keys (YubiKey)
  - JavaCard


<small>&copy;2015-2020. Green Screens Ltd. www.greenscreens.io</small>
