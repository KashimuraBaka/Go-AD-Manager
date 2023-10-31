package pshell

import "fmt"

type ADControl struct {
	*PowerShell
}

func (ps *ADControl) EnableUser(user string, unlock bool) {
	if unlock {
		ps.Execute(
			fmt.Sprintf("Enable-ADAccount '%s';", user),
			fmt.Sprintf("Get-ADUser '%s' | Move-ADObject -TargetPath 'OU=阿里,OU=催收,DC=rgrr,DC=cn';", user),
		)
	} else {
		ps.Execute(
			fmt.Sprintf("Disable-ADAccount '%s';", user),
			fmt.Sprintf("Get-ADUser '%s' | Move-ADObject -TargetPath 'OU=Disable,OU=催收,DC=rgrr,DC=cn';", user),
		)
	}
}

func (ps *ADControl) UnlockUser(user string, pwd string) {
	ps.Execute(
		fmt.Sprintf("Unlock-ADAccount -Identity '%s';", user),
		fmt.Sprintf("Set-ADAccountPassword -Identity '%s' -Reset -NewPassword (ConvertTo-SecureString -AsPlainText '%s' -Force)", user, pwd),
	)
}

func (ps *ADControl) GetUsers(user string) {
	var stdout string
	var stderr string
	var err error
	if user == "" {
		stdout, stderr, err = ps.Execute(`
		function empty {param($val, $default = $null); return $(if($val -ne $null) {$val} else {$default})};
		$ADUser = Get-ADUser -Filter * -SearchBase 'OU=催收,DC=rgrr,DC=cn' -properties logoncount,CanonicalName,lastlogon,badpasswordtime,badPwdCount,pwdLastSet;
		$ADUser = $ADUser | select name,enabled,logoncount,CanonicalName,lastlogon,badpasswordtime,badPwdCount,pwdLastSet;
		$ADComputer = Get-ADComputer -Filter * -SearchBase 'OU=催收,DC=rgrr,DC=cn' -properties Name,Description,OperatingSystem,OperatingSystemVersion;
		$ADComputer = $ADComputer | select Name,Description,OperatingSystem,OperatingSystemVersion;
		foreach($data in $ADUser) {
			$r = $ADComputer | Where-Object -FilterScript { ($_.Name -eq $data.name) } | select Description,OperatingSystem,OperatingSystemVersion;
			$netinfo = empty $r.Description '{}' | ConvertFrom-Json;
			$r.PSObject.Properties.Remove('Description');
			$r | Add-Member -Name 'ip' -Value (empty $netinfo.ip) -MemberType NoteProperty;
			$r | Add-Member -Name 'mac' -Value (empty $netinfo.mac) -MemberType NoteProperty;
			$data | Add-Member -Name 'info' -Value $r -MemberType NoteProperty;
		};
		$ADUser | ConvertTo-Json -Compress`)
	} else {
		stdout, stderr, err = ps.Execute(
			"function empty {param($val, $default = $null); return $(if($val -ne $null) {$val} else {$default})};",
			fmt.Sprintf("$ADUser = Get-ADUser '%s' -properties * | select name,enabled,logoncount,CanonicalName,lastlogon,badpasswordtime,badPwdCount,pwdLastSet;", user),
			fmt.Sprintf("$ADComputer = Get-ADComputer '%s' -properties * | select Name,Description,OperatingSystem,OperatingSystemVersion;", user),
			`$netinfo = empty $ADComputer.Description '{}' | ConvertFrom-Json;
			$ADComputer | Add-Member -Name 'ip' -Value (empty $netinfo.ip) -MemberType NoteProperty;
			$ADComputer | Add-Member -Name 'mac' -Value (empty $netinfo.mac) -MemberType NoteProperty;
			$ADComputer.PSObject.Properties.Remove('PSShowComputerName');
			$ADComputer.PSObject.Properties.Remove('RunspaceId');
			$ADComputer.PSObject.Properties.Remove('Name');
			$ADComputer.PSObject.Properties.Remove('Description');
			$ADUser | Add-Member -Name 'Info' -Value $ADComputer -MemberType NoteProperty;
			$ADUser | ConvertTo-Json -Compress;`,
		)
	}
	if err != nil {
		fmt.Println("[ADControl] Error", err)
	}
	fmt.Println(stdout)
	fmt.Println(stderr)
}

func CreateADControl() (*ADControl, error) {
	ps, err := CreatePowershell()
	if err != nil {
		return nil, err
	}

	return &ADControl{
		PowerShell: ps,
	}, nil
}
