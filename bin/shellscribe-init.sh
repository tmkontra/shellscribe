# this is run before every command
_shellscribe_trap() {
    # get last submitted command
    cmd="${history[$HISTCMD]}";
    if [ "$cmd" ]; then
        # catch the unwrap command and remove the trap
        if [ "$cmd" = "shellscribe-off" ]; then
            print 'Deactivating shellscribe';
            trap - DEBUG;
            return;
        fi
        # make sure multiple trap triggers only handle this once
        if [ "$handled" != "$HISTCMD;$cmd" ]; then
            # when either history index or command text
            # changes, we can assume its a new command
            handled="$HISTCMD;$cmd";
            # do whatever with $cmd
            outfile=$(shellscribe setup "$cmd");
	        eval $cmd | tee $outfile
        fi
        # optionally skip the raw execution
        setopt ERR_EXIT;
    fi
}

# start the debug trap
shellscribe-on() {
    print 'Activating shellscribe';
    trap '_shellscribe_trap' DEBUG;
}

# this is just defined in order to avoid errors
# the unwrapping happens in the trap handler
shellscribe-off() {}