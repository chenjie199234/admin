import { createApp } from 'vue'
import main from './main.vue'
/*
import { createVuestic } from "vuestic-ui"
// import 'vuestic-ui/css'
// createApp(main).use(createVuestic()).mount('#app')
*/
import {
	createVuesticEssential,
	VaInnerLoading,
	VaSplit,
	VaModal,
	VaModalPlugin,
	VaCard,
	VaCardTitle,
	VaCardContent,
	VaHover,
	VaButton,
	VaBadge,
	VaSelect,
	VaInput,
	VaRadio,
	VaSwitch,
	VaImage,
	VaIcon,
	VaDivider,
	VaPagination,
	VaDropdown,
	VaDropdownContent,
	VaPopover,
	VaDropdownPlugin,
} from 'vuestic-ui'
import "vuestic-ui/styles/index.css"
// import "vuestic-ui/styles/css-variables.css"
// import "vuestic-ui/styles/essential.css"
// import "vuestic-ui/styles/grid.css"
// import "vuestic-ui/styles/index2.css"
// import "vuestic-ui/styles/reset.css"
// import "vuestic-ui/styles/smart-helpers.css"
// import "vuestic-ui/styles/theme.css"
// import "vuestic-ui/styles/typography.css"

createApp(main).use(createVuesticEssential({
		components: { VaInnerLoading, VaSplit, VaModal, VaCard, VaCardTitle, VaCardContent, VaHover, VaButton, VaBadge, VaSelect, VaInput, VaRadio, VaSwitch, VaImage, VaIcon, VaDivider, VaPagination, VaDropdown, VaDropdownContent, VaPopover },
		plugins: { VaModalPlugin, VaDropdownPlugin },
		config: {},
})).mount('#app')
