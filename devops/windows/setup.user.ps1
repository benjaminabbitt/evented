function AddTo-Path
{
    param (
        [Parameter(Mandatory = $true)][String]$PathToAdd,
        [Parameter(Mandatory = $true)][ValidateSet('System', 'User')][string]$UserType,
        [Parameter(Mandatory = $true)][ValidateSet('Path', 'PSModulePath')][string]$PathType
    )

    # AddTo-Path "C:\XXX" "PSModulePath" 'System'
    if ($UserType -eq "System")
    {
        $RegPropertyLocation = 'HKLM:\System\CurrentControlSet\Control\Session Manager\Environment'
    }
    if ($UserType -eq "User")
    {
        $RegPropertyLocation = 'HKCU:\Environment'
    } # also note: Registry::HKEY_LOCAL_MACHINE\ format
    $PathOld = (Get-ItemProperty -Path $RegPropertyLocation -Name $PathType).$PathType
    "`n$UserType $PathType Before:`n$PathOld`n"
    $PathArray = $PathOld -Split ";" -replace "\\+$", ""
    if ($PathArray -notcontains $PathToAdd)
    {
        "$UserType $PathType Now:"   # ; sleep -Milliseconds 100   # Might need pause to prevent text being after Path output(!)
        $PathNew = "$PathOld;$PathToAdd"
        Set-ItemProperty -Path $RegPropertyLocation -Name $PathType -Value $PathNew
        Get-ItemProperty -Path $RegPropertyLocation -Name $PathType | select -ExpandProperty $PathType
        if ($PathType -eq "Path")
        {
            $env:Path += ";$PathToAdd"
        }                  # Add to Path also for this current session
        if ($PathType -eq "PSModulePath")
        {
            $env:PSModulePath += ";$PathToAdd"
        }  # Add to PSModulePath also for this current session
        "`n$PathToAdd has been added to the $UserType $PathType"
    }
    else
    {
        "'$PathToAdd' is already in the $UserType $PathType. Nothing to do."
    }
}

pip install pipenv --user

$packages = python -m site --user-site
$scripts = $packages.Replace("site-packages", "Scripts")
AddTo-Path "$scripts" "User" "Path"

