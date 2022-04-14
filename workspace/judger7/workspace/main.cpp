#include<bits/stdc++.h>
using namespace std;
int main()
{
    srand(time(0));
    int a=1;
    for(int i=0;i<1000000000;i++)
    {
        a=(a*(rand()%15))%10000505055;
    }
    return 0;
}