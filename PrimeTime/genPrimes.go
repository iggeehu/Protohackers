package main

func GenPrimes() map[int]bool {
		//create a map of all non prime numbers of up to a million
		//gather all values that are not part of the map
		var primeMap = make(map[int]bool)
		primeMap[1]=true
		primeMap[2]=true
		primeMap[3]=true
		var allNumbers = make(map[int]bool)
		for	i:=2; i<1000000; i++{
			//assume all numbers are true
			allNumbers[i]=true;
		}
		for i:=2; i*i<1000000; i++{
			if(allNumbers[i]==true){
				for j:=i*i; j<1000000; j+=i{
					allNumbers[j]=false;
				}
			}
		}
		for i:=2; i<1000000; i++{
			if(allNumbers[i]==true){
				primeMap[i]=true
			}
		}
		return primeMap
	}