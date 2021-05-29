# Nori Execution Steps


1. Find plugin file in directory
2. Load it and lookup for a Plugin variable
3. Cast to Plugin interface
3.1. On failure: goto step 1
4. Cast to Installable interface
4.1. On failure: go to step 5
4.2. On success: extract plugin installation record
4.2.1. If plugin not installed, add to list of installable plugins
5. 







Plugin.ID -> file
File -> Plugin.ID